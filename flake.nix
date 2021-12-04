{
  description = "Nix flake for github-hook-receiver";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";

  outputs = { self, nixpkgs }: let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in {
      defaultPackage.x86_64-linux =
        with pkgs; buildGoModule {
          pname = "github-hook-receiver";
          version = "0.0.0";
          src = ./.;

          vendorSha256 = "sha256-3aiZKr3cw8G2V3PXssewhZglAZq3f/DnzYLmb8TyLLQ=";

          subPackages = [
            "./cmd/github-hook-receiver"
          ];
        };
      devShell.x86_64-linux = with pkgs; mkShell {
        buildInputs = [
          go
        ];
      };
    };
}