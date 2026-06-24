package models

import "time"

type EstadoServicio string

const (
	EstadoProceso   EstadoServicio = "En proceso"
	EstadoEntregado EstadoServicio = "Entregado"
)

type Servicios struct {
	ID           string         `json:"id"`
	Cliente_ID   string         `json:"cliente_id"`
	TipoServicio string         `json:"tipo_servicio"`
	ModeloMarca  string         `json:"modelo_marca"`
	Precio       float64        `json:"precio"`
	Estado       EstadoServicio `json:"estado"`
	FechaPedido  time.Time      `json:"fecha_pedido"`
	// Uso un puntero ya que go no maneja NULLs nativos
	FechaEntrega  *time.Time `json:"fecha_entrega"`
	Observaciones string     `json:"observaciones"`
}

type ServiciosDetalle struct {
	ID            string  `json:"id"`
	NombreCliente string  `json:"nombre"`
	Telefono      string  `json:"telefono"`
	TipoServicio  string  `json:"tipo_servicio"`
	ModeloMarca   string  `json:"modelo_marca"`
	Precio        float64 `json:"precio"`
	Estado        string  `json:"estado"`
}
