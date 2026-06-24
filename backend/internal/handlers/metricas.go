package handlers

import (
	"Sneaker-Cleaning/internal/config"
	"encoding/json"
	"net/http"
	"time"
)

// MetricasResponse define la salida estructurada para los contadores de la app
type MetricasResponse struct {
	TotalOrdenes        int     `json:"total_ordenes"`
	IngresosDelDia      float64 `json:"ingresos_del_dia"`
	ServiciosPendientes int     `json:"servicios_pendientes"`
}

// ObtenerMetricas entrega los contadores globales requeridos por las tarjetas del frontend
func ObtenerMetricas(w http.ResponseWriter, r *http.Request) {
	// Validar que sea un método GET
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Metodo no permitido"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var metricas MetricasResponse

	// 1. Obtener Total de Órdenes históricas
	err := config.DB.QueryRow("SELECT COUNT(*) FROM servicios").Scan(&metricas.TotalOrdenes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al calcular total de ordenes: ` + err.Error() + `"}`))
		return
	}

	// 2. Obtener Servicios Pendientes o En Proceso
	err = config.DB.QueryRow("SELECT COUNT(*) FROM servicios WHERE LOWER(estado) IN ('en proceso', 'pendiente')").Scan(&metricas.ServiciosPendientes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al calcular servicios pendientes: ` + err.Error() + `"}`))
		return
	}

	// 3. Obtener Ingresos del Día (filtrado por la fecha actual en la columna fecha_pedido)
	fechaHoy := time.Now().Format("2006-01-02")
	err = config.DB.QueryRow("SELECT COALESCE(SUM(precio), 0) FROM servicios WHERE DATE(fecha_pedido) = $1", fechaHoy).Scan(&metricas.IngresosDelDia)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Error al calcular ingresos del dia: ` + err.Error() + `"}`))
		return
	}

	// Respuesta exitosa enviando el JSON estructurado al TPV
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricas)
}
