{
  description = "Personal Website for Conner Ohnesorge";

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

    bun2nix.url = "github:baileyluTCD/bun2nix";

    systems.url = "github:nix-systems/default";
  };

  nixConfig = {
    extra-substituters = ''
      https://cache.nixos.org
      https://nix-community.cachix.org
      https://devenv.cachix.org
      https://conneroisu.cachix.org
    '';
    extra-trusted-public-keys = ''
      cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY=
      nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs=
      devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=
      conneroisu.cachix.org-1:PgOlJ8/5i/XBz2HhKZIYBSxNiyzalr1B/63T74lRcU0=
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
      buildWithSpecificGo = pkg: pkg.override {inherit buildGoModule;};

      bunDeps = pkgs.callPackage ./bun.nix {};

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

            ${pkgs.bun}/bin/bun build \
                $REPO_ROOT/index.js \
                --minify \
                --minify-syntax \
                --minify-whitespace  \
                --minify-identifiers \
                --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &

            ${pkgs.tailwindcss}/bin/tailwindcss \
                --minify \
                -i $REPO_ROOT/input.css \
                -o $REPO_ROOT/cmd/conneroh/_static/dist/style.css \
                --cwd $REPO_ROOT &

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

          export PLAYWRIGHT_BROWSERS_PATH=${pkgs.playwright-driver.browsers}
          export PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
          export PLAYWRIGHT_NODEJS_PATH=${pkgs.nodejs_20}/bin/node

          # Browser executable paths
          export PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/chromium-1155"}
          export PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/firefox-1471"}
          export PLAYWRIGHT_WEBKIT_EXECUTABLE_PATH=${"${pkgs.playwright-driver.browsers}/webkit-2123"}

          echo "Playwright configured with:"
          echo "  - Browsers directory: $PLAYWRIGHT_BROWSERS_PATH"
          echo "  - Node.js path: $PLAYWRIGHT_NODEJS_PATH"
          echo "  - Chromium path: $PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH"
          echo "  - Firefox path: $PLAYWRIGHT_FIREFOX_EXECUTABLE_PATH"
          echo "  - WebKit path: $PLAYWRIGHT_WEBKIT_EXECUTABLE_PATH"

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
            inputs.bun2nix.defaultPackage.${pkgs.system}.bin

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
            bun
            nodePackages.typescript-language-server
            nodePackages.prettier

            # Infra
            flyctl
            wireguard-tools
            openssl.dev
            skopeo

            # Playwright

            playwright-driver # Provides browser archives and driver scripts
            chromium # Chromium browser
            firefox # Firefox browser
            (
              if pkgs.stdenv.isDarwin
              then pkgs.darwin.apple_sdk.frameworks.WebKit
              else pkgs.webkitgtk
            ) # WebKit browser
            nodejs_20 # Required for Playwright driver
            pkg-config # Needed for some browser dependencies
            xorg.libXcomposite # X11 Composite extension - needed by browsers
            xorg.libXdamage # X11 Damage extension - needed by browsers
            xorg.libXfixes # X11 Fixes extension - needed by browsers
            xorg.libXrandr # X11 RandR extension - needed by browsers
            xorg.libX11 # X11 client-side library
            xorg.libxcb # X11 C Bindings library
            alsa-lib # Audio library
            at-spi2-core # Accessibility support
            cairo # 2D graphics library
            cups # Printing system
            dbus # Message bus system
            expat # XML parser
            ffmpeg # Media processing
            fontconfig # Font configuration and customization
            freetype # Font rendering engine
            gdk-pixbuf # Image loading library
            glib # Low-level core library
            gtk3 # GUI toolkit
            mesa # OpenGL implementation
            nss # Network Security Services
            nspr # NetScape Portable Runtime
            pango # Text layout and rendering
          ]
          # Add the generated script packages
          ++ scriptPackages;
      };

      packages = let
        app-name = "conneroh.com";
      in rec {
        conneroh = buildGoModule {
          pname = app-name;
          name = app-name;
          version = "0.0.1";
          src = ./.;
          subPackages = ["."];
          nativeBuildInputs = [pkgs.bun];
          vendorHash = "sha256-CnE4KrZTgnUqKoB7NRPp/L+lEePlKRIx7Y/m24YzMFQ=";
          preBuild = ''
            mkdir -p node_modules
            ln -sf ${bunDeps.nodeModules}/node_modules/* node_modules/ || true
            ${scripts.nix-generate-all.exec}
          '';
        };
        C-conneroh = pkgs.dockerTools.buildLayeredImage {
          name = app-name;
          tag = "latest";
          created = "now";
          contents = [
            conneroh
            pkgs.cacert
          ];
          config = {
            WorkingDir = "/root";
            Cmd = ["/bin/conneroh.com"];
            ExposedPorts = {
              "8080/tcp" = {};
            };
            Env = [
              "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
              "NIX_SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
            ];
          };
          extraCommands = ''
            echo "$(git rev-parse HEAD)" > REVISION
          '';
        };
        deployPackage = pkgs.writeShellScriptBin "deploy" ''
          set -e

          echo "Copying image to Fly.io registry..."
          if [ -z "$FLY_AUTH_TOKEN" ]; then
            echo "FLY_AUTH_TOKEN is not set. Getting it from doppler..."
            FLY_AUTH_TOKEN=$(doppler secrets get --plain FLY_AUTH_TOKEN)
          fi

          echo "Copying image to Fly.io registry..."
          ${pkgs.skopeo}/bin/skopeo copy \
            --insecure-policy \
            docker-archive:"${C-conneroh}" \
            docker://registry.fly.io/conneroh-com:latest \
            --dest-creds x:"$FLY_AUTH_TOKEN" \
            --format v2s2

          echo "Deploying to Fly.io..."
          ${pkgs.flyctl}/bin/fly deploy \
            --remote-only \
            -c ${./fly.toml} \
            -i registry.fly.io/conneroh-com \
            -t "$FLY_AUTH_TOKEN"
        '';
      };
    });
}
