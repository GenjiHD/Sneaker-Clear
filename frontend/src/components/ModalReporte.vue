<template>
  <div class="modal-overlay-reporte">
    <div class="modal-ventana-reporte">
      <div class="reporte-header">
        <h3>📊 Historial y Reportes Globales</h3>
        <button @click="emit('cerrar')" class="btn-cerrar-X">×</button>
      </div>
      
      <p class="modal-subtitulo">Consulta todas las órdenes y filtra por rangos de fecha.</p>

      <div class="filtros-contenedor">
        <div class="filtro-campo">
          <label>Desde:</label>
          <input 
            type="date" 
            v-model="fechaDesde" 
            :max="fechaHasta" 
            @change="cargarHistorial" 
          />
        </div>
        <div class="filtro-campo">
          <label>Hasta:</label>
          <input 
            type="date" 
            v-model="fechaHasta" 
            :min="fechaDesde" 
            @change="cargarHistorial" 
          />
        </div>
        <button @click="limpiarFiltros" class="btn-limpiar" :disabled="!fechaDesde && !fechaHasta">
          🧹 Limpiar
        </button>
        
        <button 
          @click="exportarAExcel" 
          class="btn-excel" 
          :disabled="historialEntregados.length === 0"
        >
          🟢 Descargar Excel
        </button>
      </div>

      <div v-if="cargando" class="estado-reporte">Cargando historial...</div>
      
      <div v-else-if="historialEntregados.length === 0" class="estado-reporte">
        No se encontraron órdenes entregadas en este rango de fechas.
      </div>

      <div v-else class="tabla-scroll">
        <table class="tabla-reporte">
          <thead>
            <tr>
              <th>Cliente</th>
              <th>Servicio</th>
              <th>Modelo/Marca</th>
              <th>Precio</th>
              <th>Estado</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in historialEntregados" :key="item.id || item.ID">
              <td>{{ item.nombreCliente ?? item.nombre ?? 'Sin nombre' }}</td>
              <td>{{ item.tipoServicio ?? item.tipo_servicio }}</td>
              <td>{{ item.modeloMarca ?? item.modelo_marca }}</td>
              <td class="txt-precio">${{ item.precio }}</td>
              <td>
                <span class="tag-estado entregado">
                  {{ item.estado }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="resumen-footer" v-if="historialEntregados.length > 0">
        <p>Total de entregados: <strong>{{ historialEntregados.length }}</strong></p>
        <p>Monto Acumulado: <strong class="monto-total">${{ totalAcumulado }}</strong></p>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'

const emit = defineEmits(['cerrar'])
const urlAPI = import.meta.env.VITE_API_URL

const historial = ref([])
const cargando = ref(false)

const fechaDesde = ref('')
const fechaHasta = ref('')

const cargarHistorial = async () => {
  if (fechaDesde.value && fechaHasta.value && fechaHasta.value < fechaDesde.value) {
    fechaHasta.value = fechaDesde.value
  }

  cargando.value = true
  try {
    let url = `${urlAPI}/servicios/obtener`
    const params = new URLSearchParams()
    
    if (fechaDesde.value) params.append('desde', fechaDesde.value)
    if (fechaHasta.value) params.append('hasta', fechaHasta.value)
    
    if (params.toString()) {
      url += `?${params.toString()}`
    }

    const respuesta = await fetch(url)
    if (respuesta.ok) {
      historial.value = await respuesta.json()
    } else {
      console.error('Error al consultar el historial en Go')
    }
  } catch (error) {
    console.error('Error de conexión con el backend:', error)
  } finally {
    cargando.value = false
  }
}

const historialEntregados = computed(() => {
  return historial.value.filter(item => {
    const estado = (item.estado || item.Estado || '').toLowerCase()
    return estado === 'entregado'
  })
})

const limpiarFiltros = () => {
  fechaDesde.value = ''
  fechaHasta.value = ''
  cargarHistorial()
}

const totalAcumulado = computed(() => {
  return historialEntregados.value.reduce((suma, item) => suma + (item.precio || 0), 0).toFixed(2)
})

const exportarAExcel = () => {
  if (historialEntregados.value.length === 0) return

  const cabeceras = ['Cliente', 'Tipo de Servicio', 'Modelo_Marca', 'Precio', 'Estado']
  const filas = historialEntregados.value.map(item => [
    item.nombreCliente ?? item.nombre ?? 'Sin nombre',
    item.tipoServicio ?? item.tipo_servicio,
    item.modeloMarca ?? item.modelo_marca,
    item.precio,
    item.estado
  ])

  const contenidoCsv = '\uFEFF' + [
    cabeceras.join(','),
    ...filas.map(e => e.map(val => `"${String(val).replace(/"/g, '""')}"`).join(','))
  ].join('\n')

  const blob = new Blob([contenidoCsv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  
  let nombreArchivo = 'reporte_entregados'
  if (fechaDesde.value) nombreArchivo += `_desde_${fechaDesde.value}`
  if (fechaHasta.value) nombreArchivo += `_hasta_${fechaHasta.value}`
  nombreArchivo += '.csv'

  link.setAttribute('href', url)
  link.setAttribute('download', nombreArchivo)
  link.style.visibility = 'hidden'
  
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  cargarHistorial()
})
</script>

<style scoped>
/* 📌 Capa trasera oscura */
.modal-overlay-reporte {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

/* 🏠 Ventana del Modal con Altura Mejorada */
.modal-ventana-reporte {
  background-color: #ffffff;
  padding: 30px;
  border-radius: 16px;
  width: 95%;
  max-width: 850px;
  /* 📐 Ajuste de altura dinámica y cómoda */
  height: auto;
  max-height: 85vh; 
  display: flex;
  flex-direction: column;
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.4);
}

.reporte-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.reporte-header h3 {
  margin: 0;
  font-size: 22px;
  color: #0d47a1;
}

.btn-cerrar-X {
  background: none;
  border: none;
  font-size: 32px;
  cursor: pointer;
  color: #666;
  line-height: 1;
}

.modal-subtitulo {
  font-size: 14px;
  color: #555;
  margin-top: 0;
  margin-bottom: 20px;
}

/* 🔍 Contenedor de Filtros */
.filtros-contenedor {
  display: flex;
  gap: 15px;
  align-items: flex-end;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filtro-campo {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.filtro-campo label {
  font-size: 12px;
  font-weight: bold;
  color: #333;
}

.filtro-campo input {
  padding: 10px 12px;
  border: 1px solid #ccc;
  border-radius: 8px;
  background-color: #fff;
  font-size: 14px;
  color: #000000 !important;
}

.btn-limpiar {
  padding: 10px 18px;
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ccc;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
  font-size: 14px;
}

.btn-limpiar:disabled {
  background-color: #e0e0e0;
  color: #999;
  cursor: not-allowed;
}

.btn-excel {
  padding: 10px 18px;
  background-color: #1f7a42;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: bold;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-excel:hover:not(:disabled) {
  background-color: #14552c;
}

.btn-excel:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

/* 📜 SECCIÓN DE TABLA CON SCROLL VERTICAL PROPIO */
.tabla-scroll {
  flex: 1;
  overflow-y: auto; /* Activa el scroll vertical solo si la tabla se desborda */
  margin-bottom: 20px;
  border: 1px solid #eee;
  border-radius: 8px;
}

/* 📊 Estilos de la Tabla */
.tabla-reporte {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.tabla-reporte th {
  background-color: #f8f9fa;
  padding: 12px;
  font-weight: bold;
  color: #333;
  border-bottom: 2px solid #dee2e6;
  position: sticky; /* Deja las cabeceras fijas mientras haces scroll */
  top: 0;
  z-index: 10;
}

.tabla-reporte td {
  padding: 12px;
  border-bottom: 1px solid #eee;
  color: #222;
  font-size: 15px;
}

.txt-precio {
  font-weight: bold;
  color: #2e7d32;
}

.tag-estado {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
  display: inline-block;
}

.tag-estado.entregado {
  background-color: #d4edda;
  color: #155724;
}

.estado-reporte {
  text-align: center;
  padding: 4px;
  color: #666;
  font-weight: 500;
}

/* 💰 Footer Fijo abajo */
.resumen-footer {
  display: flex;
  justify-content: space-between;
  background-color: #f1f8ff;
  padding: 15px 20px;
  border-radius: 10px;
  border: 1px solid #b3d7ff;
}

.resumen-footer p {
  margin: 0;
  font-size: 16px;
  color: #0d47a1;
}

.monto-total {
  font-size: 18px;
  color: #1b5e20;
}
</style>
