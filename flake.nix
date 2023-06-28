{
  description = "OpenProject CLI";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  inputs.devshell.url = "github:numtide/devshell";
  inputs.devshell.inputs.nixpkgs.follows = "nixpkgs";
  inputs.devshell.inputs.systems.follows = "systems";

  outputs = {self, nixpkgs, systems, devshell }:
    let
      eachSystem = nixpkgs.lib.genAttrs (import systems);
      # Nixpkgs instantiated for system types in nix-systems
      nixpkgsFor = eachSystem (system:
        import nixpkgs {
          inherit system;
          overlays = [
            devshell.overlays.default
          ];
        }
      );
    in
    {
      devShells = eachSystem (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.devshell.mkShell {
            # Add additional packages you'd like to be available in your devshell
            # PATH here
            devshell.packages = with pkgs; [
              go
            ];
          };
        });

      packages = eachSystem (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          openproject-cli = pkgs.buildGoModule {
            pname = "openproject-cli";
            version = "0.0.0";
            src = ./.;
            vendorHash = "sha256-JFvC9V0xS8SZSdLsOtpyTrFzXjYAOaPQaJHdcnJzK3s="; 
            meta = with pkgs.lib; {
              description = "Simple command-line interface to OpenProject";
              homepage = "https://github.com/opf/openproject-cli";
              licenses = licenses.gpl3;
            };
          };
        });
    };
}
