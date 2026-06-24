<template>
  <div class="contenedor-principal-pro">
    
    <div class="seccion-metricas">
      <div class="tarjeta-metrica">
        <span class="metrica-titulo">Órdenes Totales</span>
        <span class="metrica-valor">{{ totalOrdenes }}</span>
      </div>
      <div class="tarjeta-metrica">
        <span class="metrica-titulo">Ingresos del Día</span>
        <span class="metrica-valor">${{ ingresosDelDia }}</span>
      </div>
      <div class="tarjeta-metrica">
        <span class="metrica-titulo">Servicios Pendientes</span>
        <span class="metrica-valor">{{ serviciosPendientes }}</span>
      </div>
    </div>

    <div class="tarjeta-tabla-pro">
      <div class="header-tabla-pro">
        <h2>Órdenes Registradas</h2>
        <button @click="mostrarModalReporte = true" class="btn-reporte-pro">
          📊 Historial y Reportes
        </button>
      </div>

      <div class="contenedor-buscador">
        <input 
          type="text" 
          v-model="filtroBusqueda" 
          placeholder="Buscar cliente o modelo..." 
          class="input-buscador"
        />
      </div>

      <div v-if="cargando" class="mensaje-estado-pro">Cargando órdenes...</div>
      <div v-else-if="ordenesFiltradas.length === 0" class="mensaje-estado-pro">
        No se encontraron órdenes.
      </div>

      <div class="tabla-responsiva" v-else>
        <table class="tabla-pro">
          <thead>
            <tr>
              <th>Folio</th>
              <th>Cliente</th>
              <th>Servicio</th>
              <th>Estado</th>
              <th>Total</th>
              <th class="texto-centrado">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(orden, index) in ordenesFiltradas" :key="orden.id">
              <td class="col-folio">#{{ orden.id.toString().slice(-4).toUpperCase() }}</td>
              <td>
                <div class="datos-cliente">
                  <span class="nombre-cliente">{{ orden.nombre }}</span>
                  <span class="tel-cliente">📞 {{ orden.telefono }}</span>
                </div>
              </td>
              <td>
                <div class="datos-servicio">
                  <span class="modelo-calzado">{{ orden.modelo_marca }}</span>
                  <span class="tipo-lavado">{{ orden.tipo_servicio }}</span>
                </div>
              </td>
              <td>
                <span class="badge-estado" :class="orden.estado.toLowerCase().replace(' ', '-')">
                  {{ orden.estado }}
                </span>
              </td>
              <td class="col-precio">${{ orden.precio }}</td>
              <td>
                <div class="acciones-celda">
                  <button @click="abrirEditor(orden)" class="btn-tabla btn-editar" title="Editar">✏️</button>
                  <button @click="verTicketPDF(orden.id)" class="btn-tabla btn-ticket" title="Imprimir Ticket">🖨️</button>
                  <button @click="eliminarOrden(orden.id)" class="btn-tabla btn-entregar" title="Marcar como Entregado">✅</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
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
const filtroBusqueda = ref('') // <-- Filtro reactivo del input
const urlAPI = import.meta.env.VITE_API_URL

const estadoModal = ref(false)
const ordenSeleccionada = ref(null)
const mostrarModalReporte = ref(false)
const mostrarModalTicket = ref(false)
const urlTicketPDF = ref('')

const obtenerOrdenes = async () => {
alert("El .exe está intentando conectar a: " + urlAPI);
  try {
    const respuesta = await fetch(`${urlAPI}/api/servicios/obtener`)
    if (respuesta.ok) {
      ordenes.value = await respuesta.json()
    }
  } catch (error) {
    console.error('Error de conexión:', error)
  } finally {
    cargando.value = false
  }
}

// 📊 CALCULOS MATEMÁTICOS DE LAS TARJETAS SUPERIORES
const totalOrdenes = computed(() => ordenes.value.length)

const ingresosDelDia = computed(() => {
  return ordenes.value
    .reduce((sum, orden) => sum + (parseFloat(orden.precio) || 0), 0)
})

const serviciosPendientes = computed(() => {
  return ordenes.value.filter(orden => {
    const est = (orden.estado || '').toLowerCase()
    return est === 'en proceso' || est === 'pendiente'
  }).length
})

// 🔍 BUSCADOR DINÁMICO
const ordenesFiltradas = computed(() => {
  return ordenes.value.filter(orden => {
    const busqueda = filtroBusqueda.value.toLowerCase()
    const nombreMatch = (orden.nombre || '').toLowerCase().includes(busqueda)
    const modeloMatch = (orden.modelo_marca || '').toLowerCase().includes(busqueda)
    // Mostramos solo lo que está activo/en proceso en la tabla principal
    const estaActivo = (orden.estado || '').toLowerCase() === 'en proceso'
    
    return estaActivo && (nombreMatch || modeloMatch)
  })
})

const abrirEditor = (orden) => {
  ordenSeleccionada.value = orden  
  estadoModal.value = true         
}

const verTicketPDF = (id) => {
  urlTicketPDF.value = `${urlAPI}/api/servicios/ticket?id=${id}`
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
    const respuesta = await fetch(`${urlAPI}/api/servicios/cambiarestado`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: id })
    })

    if (respuesta.ok) {
      await obtenerOrdenes() 
    } else {
      const errTexto = await respuesta.text()
      alert(`No se pudo actualizar: ${errTexto}`)
    }
  } catch (error) {
    alert('Error de conexión al intentar comunicar con Go.')
  }
}

onMounted(() => {
  obtenerOrdenes()
  window.addEventListener('focus', obtenerOrdenes)
})
</script>

<style scoped>
.contenedor-principal-pro {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
  font-family: inherit;
}

/* 📊 GRID DE MÉTRICAS (Estilo exacto de tu captura) */
.seccion-metricas {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.tarjeta-metrica {
  background-color: #0f172a; /* Azul ultra oscuro del prototipo */
  color: #ffffff;
  padding: 20px;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.metrica-titulo {
  font-size: 13px;
  font-weight: 500;
  color: #94a3b8; /* Gris fino */
}

.metrica-valor {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.03em;
}

/* 📋 TARJETA CONTENEDORA DE LA TABLA */
.tarjeta-tabla-pro {
  background: #ffffff;
  padding: 24px;
  border-radius: 16px;
  border: 1px solid #f3f4f6;
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.02);
}

.header-tabla-pro {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-tabla-pro h2 {
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.btn-reporte-pro {
  padding: 9px 14px;
  background-color: #ffffff;
  color: #1e293b;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-reporte-pro:hover {
  background-color: #f8fafc;
  border-color: #cbd5e1;
}

/* 🔍 INPUT BUSCADOR */
.contenedor-buscador {
  margin-bottom: 18px;
}

.input-buscador {
  width: 100%;
  padding: 11px 16px;
  font-size: 14px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background-color: #ffffff;
  color: #0f172a;
  box-sizing: border-box;
}

.input-buscador:focus {
  outline: none;
  border-color: #cbd5e1;
}

/* 📝 TABLA PROFESIONAL */
.tabla-responsiva {
  width: 100%;
  overflow-x: auto;
}

.tabla-pro {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 14px;
}

.tabla-pro th {
  padding: 12px 16px;
  font-weight: 600;
  color: #475569;
  border-bottom: 1px solid #f1f5f9;
  background-color: #f8fafc;
  font-size: 13px;
}

.tabla-pro td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #0f172a;
  vertical-align: middle;
}

/* Sub-estilos por celdas */
.col-folio {
  font-family: monospace;
  font-weight: 700;
  color: #64748b;
}

.datos-cliente, .datos-servicio {
  display: flex;
  flex-direction: column;
}

.nombre-cliente {
  font-weight: 600;
  color: #0f172a;
}

.tel-cliente {
  font-size: 12px;
  color: #64748b;
  margin-top: 2px;
}

.modelo_calzado {
  font-weight: 500;
}

.tipo-lavado {
  font-size: 12px;
  color: #0d47a1;
  font-weight: 600;
}

.col-precio {
  font-weight: 700;
  color: #0f172a;
}

/* BADGES DE ESTADO (Vistosos y limpios) */
.badge-estado {
  display: inline-block;
  padding: 4px 10px;
  font-size: 11px;
  font-weight: 700;
  border-radius: 20px;
  text-transform: uppercase;
}

.badge-estado.en-proceso, .badge-estado.pendiente {
  background-color: #fef3c7;
  color: #d97706;
}

/* 🛠️ ACCIONES CELDAS (Botones redondos tipo app profesional) */
.acciones-celda {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.btn-tabla {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid #e2e8f0;
  background-color: #ffffff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  transition: all 0.15s ease;
}

.btn-tabla:hover {
  transform: scale(1.08);
}

.btn-editar:hover { background-color: #f1f5f9; }
.btn-ticket:hover { background-color: #fef9c3; }
.btn-entregar:hover { background-color: #dcfce7; }

.mensaje-estado-pro {
  text-align: center;
  padding: 30px;
  color: #64748b;
  font-size: 14px;
}

/* VISOR PDF MODAL */
.modal-pdf-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(15, 23, 42, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.modal-pdf-contenido {
  background-color: #ffffff;
  width: 85%;
  max-width: 700px;
  height: 80vh;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-pdf-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background-color: #f8fafc;
  border-bottom: 1px solid #f1f5f9;
}

.btn-cerrar-modal {
  background: none;
  border: none;
  font-weight: 600;
  cursor: pointer;
  color: #ef4444;
}

.modal-pdf-cuerpo {
  flex: 1;
  background-color: #475569;
}

.visor-pdf {
  width: 100%;
  height: 100%;
  border: none;
}
</style>