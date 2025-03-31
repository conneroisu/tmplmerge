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

    mk-shell-bin.url = "github:rrbutani/nix-mk-shell-bin";

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
      # overlays = [(final: prev: {final.go = prev.go_1_24;})];
      pkgs = import inputs.nixpkgs {inherit system;};
      buildGoModule = pkgs.buildGoModule.override {};
      buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};

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
        build = {
          exec = ''
            nix build --accept-flake-config .#packages.x86_64-linux.conneroh
          '';
          description = "Build the package";
        };
        update = {
          exec = ''
            ${pkgs.doppler}/bin/doppler run -- ${pkgs.go}/bin/go run $REPO_ROOT/cmd/update --cwd $REPO_ROOT
          '';
          description = "Update the generated go files.";
        };
        generate-reload = {
          exec = ''
            ${pkgs.templ}/bin/templ generate &
            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i ./input.css \
                -o ./cmd/conneroh/_static/dist/style.css \
                --cwd $REPO_ROOT &
            wait
          '';
          description = "Generate templ files and wait for completion";
        };

        generate-all = {
          exec = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel) # needed
            ${buildWithSpecificGo pkgs.gomarkdoc}/bin/gomarkdoc -o README.md -e .

            wait
          '';
          description = "Generate js files";
        };

        nix-generate-all = {
          exec = ''
            ${pkgs.templ}/bin/templ generate &

            ${pkgs.bun}/bin/bun build \
                ./index.ts \
                --minify \
                --minify-syntax \
                --minify-whitespace  \
                --minify-identifiers \
                --outdir ./cmd/conneroh/_static/dist/ &

            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i ./input.css \
                -o ./cmd/conneroh/_static/dist/style.css \
                --cwd . &

            wait
          '';
          description = "Generate all files in parallel";
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
        run = {
          exec = ''cd $REPO_ROOT && air'';
          description = "Run the application with air for hot reloading";
        };
      };

      # Convert scripts to packages
      scriptPackages =
        pkgs.lib.mapAttrsToList
        (name: script: pkgs.writeShellScriptBin name script.exec)
        scripts;
    in {
      devShells.default = pkgs.mkShell {
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
            (buildWithSpecificGo gopls)
            (buildWithSpecificGo templ)
            (buildWithSpecificGo golines)
            (buildWithSpecificGo golangci-lint-langserver)
            (buildWithSpecificGo gomarkdoc)
            (buildWithSpecificGo gotests)
            (buildWithSpecificGo gotools)
            (buildWithSpecificGo reftools)
            graphviz

            # Web
            tailwindcss
            tailwindcss-language-server
            nodePackages.typescript-language-server
            nodePackages.prettier

            # Infra
            wireguard-tools
            openssl.dev
          ]
          # Add the generated script packages
          ++ scriptPackages;
      };

      packages = {
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
