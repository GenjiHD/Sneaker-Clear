{ pkgs ? import <nixpkgs> {} }:

let
  # 🎨 Librerías del sistema que Tauri necesita para renderizar la interfaz gráfica
  libraries = with pkgs; [
    webkitgtk_4_1
    gtk3
    cairo
    gdk-pixbuf
    glib
    dbus
    librsvg
  ];
in
pkgs.mkShell {
  # Herramientas para el entorno (Tus herramientas + las de Tauri/Rust)
  buildInputs = with pkgs; [
    # --- Tus herramientas actuales ---
    go
    gnumake
    air
    nodejs
    
    # --- Herramientas añadidas para Tauri/Rust ---
    cargo
    rustc
    pkg-config
    libsoup_3
  ] ++ libraries;

  # Variables de entorno al entrar al shell
  shellHook = ''
    # 🚨 CRÍTICO PARA NIXOS: Mapea las rutas de las librerías para que el compilador las encuentre
    export PKG_CONFIG_PATH="${pkgs.lib.makeSearchPathOutput "dev" "lib/pkgconfig" libraries}"
    export LD_LIBRARY_PATH="${pkgs.lib.makeLibraryPath libraries}:$LD_LIBRARY_PATH"

    echo "🛸 Entorno de desarrollo para Sneaker Cleaning (Go, Vue 3 & Tauri listo)"
  '';
}
