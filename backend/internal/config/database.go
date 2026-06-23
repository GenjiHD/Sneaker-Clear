package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// Pool de conexciones global
var DB *sql.DB

func ConnectDB() error {
	// Cargamos el archivo .env con la URL
	if err := godotenv.Load(); err != nil {
		fmt.Println(" No se encontro archivo .env, usando variables del sistema")
	}

	// Leemos la URL de conexion
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return fmt.Errorf("La variable DATABASE_URL esta vacia")
	}

	// Abrimos el pool de conexiones usando pgx
	var err error
	DB, err = sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("Error al abrir la base de datos: %v", err)
	}

	// Configuramos los limites del pool de conexion
	DB.SetMaxOpenConns(10)                  // Maximo 10 conexiones activas
	DB.SetMaxIdleConns(2)                   // Mantener maximo 2 conexiones inactivas
	DB.SetConnMaxLifetime(30 * time.Minute) // Tiempo maximo de vida de las conexiones

	// Verificar conexion
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("No se pudo conectar a Supabase: %v", err)
	}

	fmt.Println("Conexion a Supabase exitosa")
	return nil
}
