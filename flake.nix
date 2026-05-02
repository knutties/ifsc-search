{
  description = "bank-search — HTTP search service for Indian bank branches";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go_1_25
            gopls
            gotools
            golangci-lint
            gh
            jq
            smithy-cli
            nodejs_22
          ];

          shellHook = ''
            echo "bank-search dev shell ($(go version | awk '{print $3}'))"
          '';
        };
      });
}
