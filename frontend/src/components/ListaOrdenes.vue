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
          <button @click="eliminarOrden(orden.id)" class="btn-accion btn-entregar">🗑️ Entregar</button>
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
  /* Ensancho el contenedor máximo para dar espacio a las columnas colaterales */
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

/* ⚡ TRANSFORMACIÓN DE LISTA A GRID RESPONSIVO */
.lista-tickets {
  display: grid;
  /* Calcula columnas automáticas de mínimo 280px para que entren de 3 a 4 por fila */
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
  justify-content: space-between; /* Empuja el bloque de acciones siempre al final de la tarjeta */
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
  margin-left: auto; /* Desplaza la etiqueta hacia el extremo derecho */
}

/* 🏁 CONTENEDOR DE BOTONES HORIZONTALES */
.ticket-acciones {
  display: flex;
  gap: 10px;
  margin-top: auto; 
}

.btn-accion {
  flex: 1; /* Divide el espacio disponible equitativamente al 50% cada uno */
  padding: 10px;
  background-color: #ffffff;
  border: 1px solid #dcdcdc;
  border-radius: 8px;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  transition: background 0.2s, color 0.2s, border-color 0.2s;
  color: #000000 !important; /* Mantiene la visibilidad del texto en cualquier esquema */
}

.btn-accion:hover {
  background-color: #0d47a1;
  color: #ffffff !important;
  border-color: #0d47a1;
}

/* Destaca la acción final con un verde sutil al pasar el puntero */
.btn-accion.btn-entregar:hover {
  background-color: #28a745;
  border-color: #28a745;
  color: #ffffff !important;
}
</style>
