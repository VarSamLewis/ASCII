{
  description = "ascii - a minimal Go project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        buildAscii =
          {
            doCheck ? false,
          }:
          pkgs.buildGoModule {
            pname = "ascii";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-v9QGoqKB/LGeAPF4NNTyI1Nmy301m43/9ljorcayums=";

            postInstall = ''
              if [ -f $out/bin/src ]; then
                mv $out/bin/src $out/bin/ascii
              fi
            '';

            inherit doCheck;
          };
      in
      {
        packages.default = buildAscii { doCheck = false; };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.ffmpeg
          ];
        };

        checks.default = buildAscii { doCheck = true; };
      }
    );
}
