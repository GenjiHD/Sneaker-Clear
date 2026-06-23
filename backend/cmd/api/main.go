package main

import (
	"Sneaker-Cleaning/internal/config"
	"Sneaker-Cleaning/internal/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Inicializar pool de conexiones
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Error al conectar la base de datos: %v", err)
	}

	// Cerrar las conexiones cuando el server este apagado
	defer config.DB.Close()

	// Configurando servermux
	mux := http.NewServeMux()

	// Ruta de prueba
	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "online", "backend": "Go", "database": "Conectada"}`))
	})

	// Rutas funcionales
	mux.HandleFunc("POST /api/clientes/crear", handlers.RegistrarCliente)
	mux.HandleFunc("POST /api/servicios/crear", handlers.RegistrarServicio)
	mux.HandleFunc("GET /api/servicios/obtener", handlers.ListarServicios)
	mux.HandleFunc("PUT /api/servicios/cambiarestado", handlers.CambiarEstado)
	mux.HandleFunc("PUT /api/servicios/actualizar", handlers.EditarServicio)
	mux.HandleFunc("GET /api/servicios/ticket", handlers.GenerarTicketPDFHandler)
	// Configuracion de CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "PUT", "DELETE", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(mux)

	// Arrancar server
	fmt.Println("Servidor escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatalf("No se pudo levantar el server: %v", err)
	}
}
