{
  description = "Vintage Story Mod Updater";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "v2.0.2";
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      defaultPackage.${system} = pkgs.buildGoModule {
        inherit version;
        pname = "VSModUpdater";
        src = self.outPath;
        vendorHash = "sha256-GD9RRUiX21aV8RpT0kN7rEcc2atOgU9ePvWBoLL6t2E=";

        ldflags = [
          "-s"
          "-w"
          "-X github.com/rafalb8/VSModUpdater/v2/internal/config.version=${version}"
        ];

        meta = with pkgs.lib; {
          description = "Vintage Story Mod Updater";
          homepage = "https://github.com/rafalb8/VSModUpdater";
          license = licenses.mit;
          maintainers = with maintainers; [ rafalb8 ];
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        env.CGO_ENABLED = "0";
        env.TAG = version;
        packages = with pkgs; [ go zip ];
      };
    };
}
