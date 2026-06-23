package handlers

import (
	"fmt"
	"net/http"
	"time"

	"Sneaker-Cleaning/internal/config"
	"Sneaker-Cleaning/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// GenerarTicketPDFHandler maneja el endpoint GET /servicios/ticket?id=XX
// GenerarTicketPDFHandler maneja el endpoint GET /api/servicios/ticket?id=UUID
func GenerarTicketPDFHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Validar método HTTP
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// 2. Obtener el ID (UUID) del query string
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Falta el ID del servicio", http.StatusBadRequest)
		return
	}

	// 3. CONSULTAR LA BASE DE DATOS USING TU PROPIO ESTILO (JOIN CON CLIENTES)
	var s models.ServiciosDetalle
	query := `
		select s.id, c.nombre, c.telefono, s.tipo_servicio, s.modelo_marca, s.precio, s.estado
		from servicios s
		join clientes c on s.cliente_id = c.id
		where s.id = $1
	`

	// Usamos config.DB tal cual lo haces en tus otros métodos
	err := config.DB.QueryRow(query, id).Scan(&s.ID, &s.NombreCliente, &s.Telefono, &s.TipoServicio, &s.ModeloMarca, &s.Precio, &s.Estado)
	if err != nil {
		http.Error(w, "No se encontro el servicio: "+err.Error(), http.StatusNotFound)
		return
	}

	// 4. DISEÑAR EL PDF CON GOFPDF (Formato compacto A6 para ticketera)
	pdf := gofpdf.New("P", "mm", "A6", "")
	pdf.AddPage()
	pdf.SetMargins(10, 10, 10)

	// 🌟 Creamos el traductor de codificación para arreglar los caracteres raros
	traductor := pdf.UnicodeTranslatorFromDescriptor("")

	// Encabezado / Identidad de la Tienda
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "SNEAKER CLEAR", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, traductor("Cuidado y Limpieza Profesional"), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 4, fmt.Sprintf("Fecha: %s", time.Now().Format("02/01/2006 15:04")), "", 1, "C", false, 0, "")

	// Línea divisoria
	pdf.Ln(4)
	pdf.CellFormat(0, 0, "", "T", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Encabezado del ticket (mostrando los últimos 8 caracteres si es un UUID largo)
	ticketNo := s.ID
	if len(ticketNo) > 8 {
		ticketNo = ticketNo[len(ticketNo)-8:]
	}
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 5, fmt.Sprintf("TICKET DE SERVICIO #%s", ticketNo), "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Datos recopilados del cliente
	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(16, 4, "Cliente:")
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, traductor(s.NombreCliente), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 8)
	pdf.Cell(16, 4, traductor("Teléfono:"))
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, s.Telefono, "", 1, "L", false, 0, "")

	pdf.Ln(2)
	pdf.CellFormat(0, 0, "", "T", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Detalles del calzado ingresado
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(0, 4, "Detalles del Calzado:", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(0, 4, traductor(fmt.Sprintf("  - Modelo: %s", s.ModeloMarca)), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, traductor(fmt.Sprintf("  - Servicio: %s", s.TipoServicio)), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, traductor(fmt.Sprintf("  - Estado: %s", s.Estado)), "", 1, "L", false, 0, "")

	pdf.Ln(3)
	pdf.CellFormat(0, 0, "", "T", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Importe total
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(40, 6, "TOTAL:")
	pdf.CellFormat(0, 6, fmt.Sprintf("$%.2f", s.Precio), "", 1, "R", false, 0, "")

	// Mensaje de cierre
	pdf.Ln(5)
	pdf.SetFont("Arial", "I", 7)
	pdf.CellFormat(0, 3, "Conserva este comprobante para retirar tus sneakers.", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 3, traductor("¡Muchas gracias por elegirnos!"), "", 1, "C", false, 0, "")

	// 5. RESPONDER CON EL FLUJO DEL PDF EN TIEMPO REAL
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=ticket.pdf")

	err = pdf.Output(w)
	if err != nil {
		http.Error(w, "Error al procesar el archivo PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
