<template>
  <div class="contenedor-principal-ordenes">
    <div class="contenedor-header-acciones">
      <button @click="mostrarModalReporte = true" class="btn-reporte">
        📊 Historial y Reportes
      </button>
    </div>

    <h2>Órdenes en Proceso</h2>

    <div v-if="cargando" class="mensaje-estado">Cargando órdenes...</div>
    <div v-if="!cargando && ordenesEnProceso.length === 0" class="mensaje-estado">No hay órdenes pendientes.</div>

    <div class="lista-tickets" v-else>
      <div v-for="orden in ordenesEnProceso" :key="orden.id" class="ticket-tarjeta">
        
        <div class="ticket-info">
          <p class="cliente-nombre">{{ orden.nombre }}</p>
          
          <p class="servicio-detalle">
            {{ orden.modelo_marca }} - {{ orden.tipo_servicio }}
          </p>
          
          <p class="icono-texto">
            <span class="icono">📞</span> {{ orden.telefono }}
          </p>
          
          <p class="icono-texto precio">
            <span class="icono">$</span> {{ orden.precio }} 
            <span class="estado-tag">{{ orden.estado }}</span>
          </p>
        </div>

        <div class="ticket-acciones">
          <button @click="abrirEditor(orden)" class="btn-accion">✏️ Editar</button>
          <button @click="verTicketPDF(orden.id)" class="btn-accion btn-ticket">🖨️ Ticket</button>
          <button @click="eliminarOrden(orden.id)" class="btn-accion btn-entregar">🗑️ Entregar</button>
        </div>

      </div>
    </div>

    <div v-if="mostrarModalTicket" class="modal-pdf-overlay" @click.self="cerrarModalTicket">
      <div class="modal-pdf-contenido">
        <div class="modal-pdf-header">
          <h3>Vista Previa del Ticket</h3>
          <button @click="cerrarModalTicket" class="btn-cerrar-modal">❌ Cerrar</button>
        </div>
        <div class="modal-pdf-cuerpo">
          <iframe :src="urlTicketPDF" class="visor-pdf"></iframe>
        </div>
      </div>
    </div>

    <ModalEditar 
      v-if="estadoModal && ordenSeleccionada" 
      :orden="ordenSeleccionada" 
      :key="ordenSeleccionada.id"
      @cerrar="estadoModal = false" 
      @actualizado="obtenerOrdenes" 
    />

    <ModalReporte 
      v-if="mostrarModalReporte" 
      @cerrar="mostrarModalReporte = false" 
    />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import ModalEditar from './EditarOrden.vue'
import ModalReporte from './ModalReporte.vue'

const ordenes = ref([])
const cargando = ref(true)
const urlAPI = import.meta.env.VITE_API_URL

const estadoModal = ref(false)
const ordenSeleccionada = ref(null)
const mostrarModalReporte = ref(false)

// Estados para el Modal del PDF
const mostrarModalTicket = ref(false)
const urlTicketPDF = ref('')

const obtenerOrdenes = async () => {
  try {
    const respuesta = await fetch(`${urlAPI}/servicios/obtener`)
    if (respuesta.ok) {
      const datos = await respuesta.json()
      ordenes.value = datos
    } else {
      console.error('Error al traer datos de Go')
    }
  } catch (error) {
    console.error('Error de conexión:', error)
  } finally {
    cargando.value = false
  }
}

// Filtra solo las órdenes que están en proceso sin importar capitalización
const ordenesEnProceso = computed(() => {
  return ordenes.value.filter(orden => {
    const estado = (orden.estado || orden.Estado || '').toLowerCase()
    return estado === 'en proceso'
  })
})

const abrirEditor = (orden) => {
  ordenSeleccionada.value = orden  
  estadoModal.value = true         
}

// Función para abrir el modal y setear la URL del PDF generado por Go
const verTicketPDF = (id) => {
  urlTicketPDF.value = `${urlAPI}/servicios/ticket?id=${id}`
  mostrarModalTicket.value = true
}

const cerrarModalTicket = () => {
  mostrarModalTicket.value = false
  urlTicketPDF.value = ''
}

const eliminarOrden = async (id) => {
  const confirmar = confirm('¿Estás seguro de marcar este servicio como Entregado?')
  if (!confirmar) return

  try {
    const respuesta = await fetch(`${urlAPI}/servicios/cambiarestado`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ id: id })
    })

    if (respuesta.ok) {
      alert('¡Servicio entregado con éxito! Desaparecerá de la lista de pendientes.')
      await obtenerOrdenes() 
    } else {
      const errTexto = await respuesta.text()
      alert(`No se pudo actualizar el estado: ${errTexto}`)
    }
  } catch (error) {
    console.error('Error al intentar actualizar estado:', error)
    alert('Error de conexión al intentar comunicar con Go.')
  }
}

onMounted(() => {
  obtenerOrdenes()
  window.addEventListener('focus', obtenerOrdenes)
})
</script>

<style scoped>
.contenedor-principal-ordenes {
  max-width: 1200px;
  margin: 0 auto;
  font-family: sans-serif;
  background-color: #ffffff; 
  padding: 25px;
  border-radius: 16px;
  border: 1px solid #eeeeee; 
}

.contenedor-header-acciones {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 15px;
}

.btn-reporte {
  padding: 10px 16px;
  background-color: #0d47a1;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-reporte:hover {
  background-color: #0a357a;
}

h2 {
  font-size: 28px;
  font-weight: bold;
  color: #000000;
  margin-top: 0;
  margin-bottom: 25px;
}

.mensaje-estado {
  text-align: center;
  color: #666;
  padding: 20px;
}

.lista-tickets {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.ticket-tarjeta {
  background-color: #f9f9f9;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 15px;
  border: 1px solid #f0f0f0;
}

.cliente-nombre {
  font-size: 18px;
  font-weight: bold;
  color: #000000;
  margin: 0 0 4px 0;
}

.servicio-detalle {
  font-size: 15px;
  color: #333333;
  margin: 0 0 6px 0;
  font-weight: 500;
}

.icono-texto {
  font-size: 14px;
  color: #555555;
  margin: 4px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.precio {
  font-size: 16px;
  font-weight: bold;
  color: #1b5e20;
}

.estado-tag {
  font-size: 11px;
  background-color: #e3f2fd;
  color: #0d47a1;
  padding: 3px 8px;
  border-radius: 12px;
  font-weight: bold;
  margin-left: auto;
}

/* CONTENEDOR DE BOTONES EN HORIZONTAL (AJUSTADO PARA 3 BOTONES) */
.ticket-acciones {
  display: flex;
  gap: 6px;
  margin-top: auto; 
}

.btn-accion {
  flex: 1; 
  padding: 10px 4px;
  background-color: #ffffff;
  border: 1px solid #dcdcdc;
  border-radius: 8px;
  font-size: 12px; /* Un pelín más pequeña para que quepan cómodos los 3 botones */
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2px;
  transition: background 0.2s, color 0.2s, border-color 0.2s;
  color: #000000 !important;
}

.btn-accion:hover {
  background-color: #0d47a1;
  color: #ffffff !important;
  border-color: #0d47a1;
}

.btn-accion.btn-ticket:hover {
  background-color: #f1c40f;
  border-color: #f1c40f;
  color: #000000 !important;
}

.btn-accion.btn-entregar:hover {
  background-color: #28a745;
  border-color: #28a745;
  color: #ffffff !important;
}

/* ESTILOS DEL NUEVO MODAL DEL VISOR PDF */
.modal-pdf-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.modal-pdf-contenido {
  background-color: #ffffff;
  width: 80%;
  max-width: 750px;
  height: 85vh;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 5px 15px rgba(0,0,0,0.3);
}

.modal-pdf-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #eee;
}

.modal-pdf-header h3 {
  margin: 0;
  font-size: 18px;
  color: #333;
}

.btn-cerrar-modal {
  background: none;
  border: none;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  color: #c0392b;
}

.modal-pdf-cuerpo {
  flex: 1;
  background-color: #525659; /* Fondo gris oscuro típico de visores de PDF */
}

.visor-pdf {
  width: 100%;
  height: 100%;
  border: none;
}
</style>