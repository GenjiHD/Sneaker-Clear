package handlers

import (
	"Sneaker-Cleaning/internal/config"
	"Sneaker-Cleaning/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegistrarServicio(w http.ResponseWriter, r *http.Request) {
	// Validar que sea un POST
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Metodo no permitido"}`))
		return
	}

	// Estructura para recibir el JSON (Adaptada con nombre y telefono que manda Vue)
	var input struct {
		Nombre       string  `json:"nombre"`
		Telefono     string  `json:"telefono"`
		TipoServicio string  `json:"tipo_servicio"`
		ModeloMarca  string  `json:"modelo_marca"`
		Precio       float64 `json:"precio"` // Cambiado a float64 para coincidir con tu struct Servicios
	}

	// Decodificacion del JSON
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"JSON mal formado"}`))
		return
	}

	// 👥 PASO 1: Buscar si el cliente existe por su teléfono, si no existe lo inserta en caliente
	// Esta consulta atómica te regresa el UUID (id) correspondiente sin duplicar datos.
	var clienteID string
	queryCliente := `
        WITH existente AS (
            SELECT id FROM clientes WHERE telefono = $1
        ), insertado AS (
            INSERT INTO clientes (nombre, telefono)
            SELECT $2, $1
            WHERE NOT EXISTS (SELECT 1 FROM existente)
            RETURNING id
        )
        SELECT id FROM existente UNION ALL SELECT id FROM insertado;
    `

	err = config.DB.QueryRow(queryCliente, input.Telefono, input.Nombre).Scan(&clienteID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al procesar el cliente: ` + err.Error() + `"}`))
		return
	}

	// 👟 PASO 2: Insertar la orden usando tu consulta original (Ligada al clienteID encontrado)
	query := "insert into servicios(cliente_id, tipo_servicio, modelo_marca, precio, estado) values ($1, $2, $3, $4, $5) returning id"

	var servicioID string
	err = config.DB.QueryRow(
		query,
		clienteID, // 👈 Aquí inyectamos de forma segura el UUID obtenido en el paso anterior
		input.TipoServicio,
		input.ModeloMarca,
		input.Precio,
		models.EstadoProceso,
	).Scan(&servicioID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al guardar el servicio: ` + err.Error() + `"}`))
		return
	}

	// Respuesta exitosa de creación
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Servicio registrado con exito",
		"servicio_id": servicioID,
	})
}

func ListarServicios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)

		return

	}

	// Tomar los parametos de la URL: &desde = YYYY-MM-DD&hasta = YYYY-MM-DD

	desde := r.URL.Query().Get("desde")

	hasta := r.URL.Query().Get("hasta")

	query := `

        select s.id, c.nombre, c.telefono, s.tipo_servicio, s.modelo_marca, s.precio, s.estado

        from servicios s

        join clientes c on s.cliente_id = c.id

        where 1=1

    `

	// Valores de los filtros ordenados

	var args []interface{}

	paramCount := 1

	if desde != "" {

		query += fmt.Sprintf(" and s.fecha_pedido >= $%d", paramCount)

		args = append(args, desde)

		paramCount++

	}

	if hasta != "" {

		query += fmt.Sprintf(" and s.fecha_pedido <= $%d", paramCount)

		args = append(args, hasta)

		paramCount++

	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {

		http.Error(w, "Error al consultar el historial: "+err.Error(), 500)

		return

	}

	defer rows.Close()

	servicios := []models.ServiciosDetalle{}

	for rows.Next() {

		var s models.ServiciosDetalle

		err := rows.Scan(&s.ID, &s.NombreCliente, &s.Telefono, &s.TipoServicio, &s.ModeloMarca, &s.Precio, &s.Estado)
		if err != nil {

			fmt.Printf("Error al scan: ", err)

			continue

		}

		servicios = append(servicios, s)

	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(servicios)
}

func CambiarEstado(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Metodo no permitido"}`))
		return
	}

	var solicitud struct {
		ID string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&solicitud)
	if err != nil || solicitud.ID == "" {
		http.Error(w, "ID de servicio requerido o formato JSON invalido", http.StatusBadRequest)
		return
	}

	query := `update servicios set estado = 'Entregado', fecha_entrega = now() where id = $1`

	result, err := config.DB.Exec(query, solicitud.ID)
	if err != nil {
		http.Error(w, "Error al actualizar el estado en la base de datos: "+err.Error(), 500)
		return
	}

	rowsUpdated, _ := result.RowsAffected()
	if rowsUpdated == 0 {
		http.Error(w, "No se encontro ningun servicio con el ID especificado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Servicio entregado con exito! fecha de entrega registrada.",
	})
}

func EditarServicio(w http.ResponseWriter, r *http.Request) {
	// Validar que sea un método PUT o POST para actualizaciones
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Metodo no permitido"}`))
		return
	}

	// Estructura para capturar los datos editados desde el frontend
	var input struct {
		ID           string  `json:"id"`
		Nombre       string  `json:"nombre"`
		Telefono     string  `json:"telefono"`
		TipoServicio string  `json:"tipo_servicio"`
		ModeloMarca  string  `json:"modelo_marca"`
		Precio       float64 `json:"precio"`
	}

	// Decodificar el JSON entrante
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"ID de servicio requerido o JSON invalido"}`))
		return
	}

	// Qeury a la base de datos
	tx, err := config.DB.Begin()
	if err != nil {
		http.Error(w, "Error al iniciar la transaccion: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Asegurar que si algo truena en el camino, se cancelen los cambios parciales
	defer tx.Rollback()

	// ACTUALIZAR LA TABLA DE CLIENTES
	// Buscamos al cliente asociado a ese servicio específico y actualizamos su nombre y teléfono
	queryCliente := `
        UPDATE clientes 
        SET nombre = $1, telefono = $2 
        WHERE id = (SELECT cliente_id FROM servicios WHERE id = $3)
    `
	_, err = tx.Exec(queryCliente, input.Nombre, input.Telefono, input.ID)
	if err != nil {
		http.Error(w, "Error al actualizar los datos del cliente: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ACTUALIZAR LA TABLA DE SERVICIOS
	// Actualizamos los detalles del calzado y el cobro
	queryServicio := `
        UPDATE servicios 
        SET tipo_servicio = $1, modelo_marca = $2, precio = $3 
        WHERE id = $4
    `
	_, err = tx.Exec(queryServicio, input.TipoServicio, input.ModeloMarca, input.Precio, input.ID)
	if err != nil {
		http.Error(w, "Error al actualizar los datos del servicio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// GUARDAR CAMBIOS DEFINITIVOS SI TODO SALIÓ BIEN
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error al confirmar la edicion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respuesta de éxito
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Orden y datos de cliente actualizados con exito.",
	})
}
