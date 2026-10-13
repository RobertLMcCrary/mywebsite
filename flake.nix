{
  description = "Portfolio site";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    devshell = {
      url = "github:numtide/devshell";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      devshell,
    }:
    let
      systems = [
        "aarch64-darwin"
        "aarch64-linux"
        "x86_64-darwin"
        "x86_64-linux"
      ];

      eachSystem =
        f:
        nixpkgs.lib.genAttrs systems (
          system:
          f {
            inherit system;
            pkgs = nixpkgs.legacyPackages.${system};
          }
        );
    in
    {
      devShells = eachSystem (
        { pkgs, system, ... }:
        {
          default = (devshell.legacyPackages.${system}.mkShell) {
            name = "portfolio";

            # Non-secret defaults — safe to commit.
            env = [
              {
                name = "ENV";
                value = "development";
              }
              {
                name = "LISTEN_ADDR";
                value = ":3000";
              }
              {
                name = "BASE_URL";
                value = "http://localhost:3000";
              }
              {
                name = "LOG_LEVEL";
                value = "debug";
              }
              {
                name = "POSTS_DIR";
                value = "./content/posts";
              }
            ];

            packages = with pkgs; [
              go_1_26
              air
              gopls
              golangci-lint
              gotools # goimports, etc.
              nixfmt
              nil
            ];

            commands = [
              {
                name = "dev";
                category = "server";
                help = "Start Go server with hot reload (air)";
                command = "air";
              }
              {
                name = "new-post";
                category = "content";
                help = "Scaffold a new blog post markdown file";
                command = ''
                  slug="$1"
                  if [ -z "$slug" ]; then
                    echo "usage: new-post <slug>"
                    exit 1
                  fi
                  date=$(date +%Y-%m-%d)
                  file="$POSTS_DIR/''${date}-''${slug}.md"
                  mkdir -p "$POSTS_DIR"
                  cat > "$file" <<EOF
                  ---
                  title: "TODO"
                  date: $date
                  ---

                  EOF
                  echo "created $file"
                '';
              }
              {
                name = "fmt";
                category = "quality";
                help = "Format Go and Nix";
                command = "go fmt ./... && nixfmt flake.nix";
              }
              {
                name = "check";
                category = "quality";
                help = "Vet, lint, and build the Go module";
                command = "go vet ./... && golangci-lint run ./... && go build ./...";
              }
              {
                name = "tidy";
                category = "quality";
                help = "Sync go.mod and go.sum";
                command = "go mod tidy";
              }
              {
                name = "build";
                category = "server";
                help = "Build the production binary";
                command = "go build -o bin/portfolio ./cmd/server";
              }
            ];
          };
        }
      );

      formatter = eachSystem ({ pkgs, ... }: pkgs.nixfmt);
    };
}
