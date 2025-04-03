// Package main is the main package for the hash command.
package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Command line flags
	dirPath         = flag.String("dir", "", "Path to the directory to check for changes")
	verbose         = flag.Bool("v", false, "Enable verbose output")
	excludePatterns = flag.String("exclude", "", "Comma-separated list of glob patterns to exclude")
	hashFilePath    = flag.String("cache", "", "Path to the cache file (defaults to .dir_hash.json in the directory)")
)

const defaultHashFileName = ".cache.json"

// Cache is the storage of previous hashes.
type Cache struct {
	HashFile string `json:"-"`
	Hashes   map[string]string
}

// Close writes the config to disk and closes the file.
func (c *Cache) Close() (err error) {
	body, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = os.WriteFile(c.HashFile, body, 0644)
	if err != nil {
		return
	}
	return
}

// LoadCache attempts to load a cache from the specified file path
func LoadCache(hashFilePath string) (*Cache, error) {
	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a new cache if the file doesn't exist
			return &Cache{
				HashFile: hashFilePath,
				Hashes:   make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	cache.HashFile = hashFilePath

	return &cache, nil
}
func main() {
	if err := run(); err != nil {
		if err == context.Canceled {
			log.Println("Operation was canceled")
			os.Exit(1)
		}
		log.Fatalf("Error: %v", err)
	}
	// Exit with code 0 to indicate no changes
	os.Exit(0)
}

func run() error {
	flag.Parse()

	// Get directory path from flag or positional argument
	dirPathValue := *dirPath
	if dirPathValue == "" {
		return fmt.Errorf("directory path is required. Use -dir flag or provide as positional argument")
	}

	// Process exclude patterns
	var excludes []string
	if *excludePatterns != "" {
		excludes = strings.Split(*excludePatterns, ",")
		if *verbose {
			fmt.Println("Excluding patterns:", excludes)
		}
	}

	// Check if directory exists
	info, err := os.Stat(dirPathValue)
	if err != nil {
		return fmt.Errorf("cannot access directory %s: %w", dirPathValue, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPathValue)
	}

	// Use default hash file path if not specified
	hashFilePathValue := *hashFilePath
	if hashFilePathValue == "" {
		hashFilePathValue = filepath.Join(dirPathValue, defaultHashFileName)
	}

	// Load or create cache
	cache, err := LoadCache(hashFilePathValue)
	if err != nil {
		return fmt.Errorf("error loading cache: %w", err)
	}

	// Calculate the hash of the directory
	currentHash, err := calculateDirectoryHash(dirPathValue, excludes)
	if err != nil {
		if err == context.Canceled {
			return err
		}
		return fmt.Errorf("error calculating directory hash: %w", err)
	}

	if *verbose {
		fmt.Printf("Current hash of %s: %s\n", dirPathValue, currentHash)
	}

	// Get the previous hash for this directory
	previousHash := cache.Hashes[dirPathValue]

	if previousHash == "" {
		if *verbose {
			fmt.Println("No previous hash found")
		}
		// First run, update the cache and exit with code 0
		cache.Hashes[dirPathValue] = currentHash
		if err := cache.Close(); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}
		fmt.Println("Initial hash created")
		os.Exit(0)
	}

	// Compare hashes
	if currentHash != previousHash {
		if *verbose {
			fmt.Printf("Changes detected in %s\n", dirPathValue)
			fmt.Printf("Previous hash: %s\n", previousHash)
			fmt.Printf("Current hash: %s\n", currentHash)
		} else {
			fmt.Printf("Changes detected in %s\n", dirPathValue)
		}

		// Update the cache
		cache.Hashes[dirPathValue] = currentHash
		if err := cache.Close(); err != nil {
			return fmt.Errorf("error writing cache: %w", err)
		}

		// Exit with code 1 to indicate changes were detected
		os.Exit(1)
	}
	fmt.Printf("No changes detected in %s\n", dirPathValue)
	return nil
}

// calculateDirectoryHash computes a hash of all files in the directory using MD5
func calculateDirectoryHash(dirPath string, excludes []string) (string, error) {
	var (
		pattern    string
		matched    bool
		relPath    string
		fileHash   string
		hasher     = md5.New()
		file       *os.File
		fileHasher = md5.New()
	)
	walkErr := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip the hash file itself
		if d.Name() == defaultHashFileName {
			return nil
		}
		// Skip directories themselves (we'll recurse into them)
		if d.IsDir() {
			return nil
		}
		// Check if this path should be excluded
		relPath, err = filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}
		for _, pattern = range excludes {
			matched, err = filepath.Match(pattern, relPath)
			if err != nil {
				return err
			}
			if matched {
				return nil
			}
		}
		// Calculate file hash
		defer fileHasher.Reset()
		fileHash, err = calculateFileHash(fileHasher, file, path)
		if err != nil {
			return err
		}
		_, err = io.WriteString(hasher, fileHash+"\n")
		if err != nil {
			return err
		}
		return nil
	})
	if walkErr != nil {
		if walkErr == context.Canceled {
			return "", walkErr
		}
		return "", walkErr
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}

// calculateFileHash computes the MD5 hash of a single file
func calculateFileHash(hasher hash.Hash, file *os.File, filePath string) (string, error) {
	var err error
	file, err = os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Printf("Error closing file: %s", err)
		}
	}()

	// Use io.Copy to efficiently copy file content to hasher
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash), nil
}
