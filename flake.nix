{
  description = "Twerge Golang Nix Flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };

    nix2container = {
      url = "github:nlewo/nix2container";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };

    systems.url = "github:nix-systems/default";
  };

  nixConfig = {
    extra-substituters = ''
      https://cache.nixos.org
      https://nix-community.cachix.org
      https://devenv.cachix.org
      https://twerge.cachix.org
    '';
    extra-trusted-public-keys = ''
      cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=
      nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs=
      devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=
      twerge.cachix.org-1:rK2EdKDH7P2S4xNTXD58XiXDpXkNr3H0rpx8huCJ9+I=
    '';
    extra-experimental-features = "nix-command flakes";
  };

  outputs = inputs @ {flake-utils, ...}:
    flake-utils.lib.eachSystem [
      "x86_64-linux"
      "i686-linux"
      "x86_64-darwin"
      "aarch64-linux"
      "aarch64-darwin"
    ] (system: let
      overlays = [(final: prev: {final.go = prev.go_1_24;})];
      pkgs = import inputs.nixpkgs {inherit system overlays;};
      buildGoModule = pkgs.buildGoModule.override {go = pkgs.go_1_24;};
      specificGo = pkg: pkg.override {inherit buildGoModule;};
    in rec {
      devShells.default = let
        scripts = {
          dx = {
            exec = ''$EDITOR $REPO_ROOT/flake.nix'';
            description = "Edit flake.nix";
          };
          clean = {
            exec = ''${pkgs.git}/bin/git clean -fdx'';
            description = "Clean Project";
          };
          tests = {
            exec = ''${pkgs.go}/bin/go test -v ./...'';
            description = "Run all go tests";
          };
          lint = {
            exec = ''
              ${pkgs.golangci-lint}/bin/golangci-lint run
              ${pkgs.statix}/bin/statix check $REPO_ROOT/flake.nix
              ${pkgs.deadnix}/bin/deadnix $REPO_ROOT/flake.nix
            '';
            description = "Run golangci-lint";
          };
          unit-tests = {
            exec = ''
              ${pkgs.go}/bin/go test -v ./...
            '';
            description = "Run unit tests.";
          };
          coverage-tests = {
            exec = ''
              ${pkgs.go}/bin/go test -v -coverprofile=coverage.out ./...
            '';
            description = "Run coverage tests.";
          };
          generate-all = {
            exec = ''
              export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
              ${specificGo pkgs.gomarkdoc}/bin/gomarkdoc -o README.md -e .
              wait
            '';
            description = "Generate js files";
          };
          format = {
            exec = ''
              cd $(git rev-parse --show-toplevel)

              ${pkgs.go}/bin/go fmt ./...

              ${pkgs.git}/bin/git ls-files \
                --others \
                --exclude-standard \
                --cached \
                -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
                | xargs prettier --write

              ${pkgs.golines}/bin/golines \
                -l \
                -w \
                --max-len=80 \
                --shorten-comments \
                --ignored-dirs=.direnv .

              cd -
            '';
            description = "Format code files";
          };
        };

        # Convert scripts to packages
        scriptPackages =
          pkgs.lib.mapAttrsToList
          (name: script: pkgs.writeShellScriptBin name script.exec)
          scripts;
      in
        pkgs.mkShell {
          shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            export CGO_CFLAGS="-O2"

            # Print available commands
            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (
                name: script: ''echo "  ${name} - ${script.description}"''
              )
              scripts
            )}
          '';
          packages = with pkgs;
            [
              # Nix
              alejandra
              nixd
              statix
              deadnix

              # Go Tools
              go_1_24
              air
              templ
              pprof
              revive
              golangci-lint
              (specificGo gopls)
              (specificGo templ)
              (specificGo golines)
              (specificGo golangci-lint-langserver)
              (specificGo gomarkdoc)
              (specificGo gotests)
              (specificGo gotools)
              (specificGo reftools)

              # Web
              tailwindcss
              tailwindcss-language-server
              nodePackages.prettier

              # Infra
              wireguard-tools
              openssl.dev
            ]
            # Add the generated script packages
            ++ scriptPackages;
        };

      overlays = {
        default = final: prev: {
          inherit (packages) hasher;
        };
      };

      packages = {
        hasher = buildGoModule {
          name = "hasher";
          src = ./cmd/hasher;
          vendorHash = null;
          version = "0.0.1";
          subPackages = ["."];
        };
        doc = pkgs.stdenv.mkDerivation {
          pname = "twerge-docs";
          version = "0.1";
          src = ./.;
          nativeBuildInputs = with pkgs; [
            nixdoc
            mdbook
            mdbook-open-on-gh
            mdbook-cmdrun
            git
          ];
          dontConfigure = true;
          dontFixup = true;
          env.RUST_BACKTRACE = 1;
          buildPhase = ''
            runHook preBuild
            cd doc  # Navigate to the doc directory during build
            mkdir -p .git  # Create .git directory
            mdbook build
            runHook postBuild
          '';
          installPhase = ''
            runHook preInstall
            mv book $out
            runHook postInstall
          '';
        };
      };
    });
}
