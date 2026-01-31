{
  description = "OpenRouter CLI - Access 400+ AI models from your terminal";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          name = "openrouter-cli-dev";
          buildInputs = with pkgs; [
            go
            gcc
            pkg-config
            git
          ];

          shellHook = ''
            echo "OpenRouter CLI development environment loaded"
            echo "Available commands: go build, make build, make install"
          '';
        };
      }
    );
}
