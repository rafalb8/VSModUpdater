{
  description = "Vintage Story Mod Updater";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "v1.0.9";
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.VSModUpdater = pkgs.buildGoModule {
        inherit version;
        pname = "VSModUpdater";
        src = ./.;
        vendorHash = null;

        buildPhase = ''
          runHook preBuild
          make linux TAG=${version}
          runHook postBuild
        '';

        installPhase = ''
          runHook preInstall
          mkdir -p $out/bin
          cp result/VSModUpdater $out/bin/VSModUpdater
          runHook postInstall
        '';

        meta = with pkgs.lib; {
          description = "Vintage Story Mod Updater";
          homepage = "https://github.com/rafalb8/VSModUpdater";
          license = licenses.mit;
          maintainers = with maintainers; [ rafalb8 ];
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        env.CGO_ENABLED = "0";
        packages = with pkgs; [ go ];
      };
    };
}
