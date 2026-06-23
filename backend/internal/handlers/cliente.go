package handlers

import (
	"encoding/json"
	"net/http"

	"Sneaker-Cleaning/internal/config"
)

func RegistrarCliente(w http.ResponseWriter, r *http.Request) {
	// Estructura temporal para recibir el json
	var input struct {
		Nombre   string `json:"nombre"`
		Telefono string `json:"telefono"`
	}

	// Decodificamos el stream de datos
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"JSON invalido o mal formado"}`))
		return
	}

	// Validacion de reglas de negocio (filtros de seguridad)
	if input.Nombre == "" || len(input.Telefono) != 10 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"El nombre es obligatorio y el telefono debe tener 10 caracteres"}`))
		return
	}

	query := "insert into clientes (nombre, telefono) values ($1, $2) returning id"

	var clienteID string

	err = config.DB.QueryRow(query, input.Nombre, input.Telefono).Scan(&clienteID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error en la base de datos: ` + err.Error() + `"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201

	respuesta := map[string]string{
		"message":    "Cliente registrado exitosamente!",
		"cliente_id": clienteID,
	}
	json.NewEncoder(w).Encode(respuesta)
}
