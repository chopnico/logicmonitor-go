with import <nixpkgs> { };
let
  packages = pkgs;
in
packages.mkShell {
  buildInputs = [
    packages.go_1_17
  ];
}
