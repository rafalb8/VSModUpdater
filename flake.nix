{
  description = "Vintage Story Mod Updater";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "v1.4.0";
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      defaultPackage.${system} = pkgs.buildGoModule {
        inherit version;
        pname = "VSModUpdater";
        src = self.outPath;
        vendorHash = "sha256-sLs/k4HovQoq6JB5jFoOFwiUmVNt7vhoGwjdvfk1fgA=";

        ldflags = [
          "-s"
          "-w"
          "-X github.com/rafalb8/VSModUpdater/internal/config.version=${version}"
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
