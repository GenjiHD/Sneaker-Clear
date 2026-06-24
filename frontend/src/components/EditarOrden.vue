<template>
  <div class="modal-overlay-edicion">
    <div class="modal-ventana-edicion">
      <h3 style="color: #0d47a1;">✏️ Modo Edición Activo</h3>
      <p class="modal-subtitulo">Modifica los campos necesarios para actualizar el ticket.</p>

      <form @submit.prevent="guardarEdicion" autocomplete="off">
        
        <div class="campo-grupo">
          <label for="input-nombre">Nombre del Cliente</label>
          <input 
            type="text" 
            id="input-nombre"
            name="ticket_nombre"
            v-model="form.nombre" 
            @input="validarNombre"
            placeholder="Ej: Carlos Pérez"
            required 
          />
        </div>

        <div class="campo-grupo">
          <label for="input-telefono">Teléfono</label>
          <input 
            type="text" 
            id="input-telefono"
            name="ticket_telefono"
            v-model="form.telefono" 
            @input="validarTelefono"
            maxlength="10"
            placeholder="Ej: 3312345678"
            required 
          />
          <small v-if="form.telefono.length > 0 && form.telefono.length < 10" class="error-texto">
            Faltan {{ 10 - form.telefono.length }} dígitos.
          </small>
        </div>

        <div class="campo-grupo">
          <label for="input-modelo">Modelo / Marca</label>
          <input 
            type="text" 
            id="input-modelo"
            name="ticket_modelo"
            v-model="form.modelo_marca" 
            required 
          />
        </div>

        <div class="campo-grupo">
          <label for="select-servicio">Tipo de Servicio</label>
          <select id="select-servicio" name="ticket_servicio" v-model="form.tipo_servicio">
            <option value="Lavado Básico">Lavado Básico</option>
            <option value="Lavado Premium">Lavado Premium</option>
            <option value="Blanqueado">Blanqueado</option>
          </select>
        </div>

        <div class="campo-grupo">
          <label for="input-precio">Precio ($)</label>
          <input 
            type="number" 
            id="input-precio"
            name="ticket_precio"
            v-model.number="form.precio" 
            required 
            min="0" 
            step="0.01" 
          />
        </div>

        <div class="modal-botones">
          <button type="button" @click="emit('cerrar')" class="btn-cancelar">Cancelar</button>
          <button type="submit" class="btn-guardar" :disabled="guardando || form.telefono.length !== 10">
            {{ guardando ? 'Guardando...' : 'Guardar Cambios' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const props = defineProps({
  orden: Object
})

const emit = defineEmits(['cerrar', 'actualizado'])

// 🚀 Inicialización directa desde la prop
const form = ref({
  id: props.orden?.id ?? props.orden?.ID ?? '',
  nombre: props.orden?.nombre ?? props.orden?.Nombre ?? props.orden?.nombreCliente ?? '',
  telefono: props.orden?.telefono ?? props.orden?.Telefono ?? '',
  modelo_marca: props.orden?.modelo_marca ?? props.orden?.Modelo_marca ?? props.orden?.ModeloMarca ?? '',
  tipo_servicio: props.orden?.tipo_servicio ?? props.orden?.Tipo_servicio ?? props.orden?.TipoService ?? props.orden?.TipoServicio ?? 'Lavado Básico',
  precio: props.orden?.precio ?? props.orden?.Precio ?? 0
})

const guardando = ref(false)
const urlAPI = import.meta.env.VITE_API_URL

// 🔥 FILTRO: Borra en tiempo real cualquier cosa que NO sea una letra (incluye acentos y ñ) o espacio
const validarNombre = () => {
  form.value.nombre = form.value.nombre.replace(/[^a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]/g, '')
}

// 🔥 FILTRO: Borra en tiempo real cualquier cosa que NO sea un número entero
const validarTelefono = () => {
  form.value.telefono = form.value.telefono.replace(/\D/g, '')
}

const guardarEdicion = async () => {
  // Doble check de seguridad por si acaso
  if (form.value.telefono.length !== 10) {
    alert('El teléfono debe tener exactamente 10 dígitos.')
    return
  }

  guardando.value = true
  try {
    const respuesta = await fetch(`${urlAPI}/api/servicios/actualizar`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })

    if (respuesta.ok) {
      alert('¡Orden actualizada con éxito!')
      emit('actualizado') 
      emit('cerrar')     
    } else {
      const errTexto = await respuesta.text()
      alert(`Error al editar: ${errTexto}`)
    }
  } catch (error) {
    console.error('Error al conectar:', error)
    alert('No se pudo conectar con el servidor de Go.')
  } finally {
    guardando.value = false
  }
}
</script>

<style scoped>
.modal-overlay-edicion {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-ventana-edicion {
  background-color: #ffffff;
  padding: 25px;
  border-radius: 16px;
  width: 90%;
  max-width: 450px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.modal-subtitulo {
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
}

.campo-grupo {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 16px;
}

.campo-grupo label {
  font-size: 14px;
  font-weight: bold;
  color: #333;
}

.campo-grupo input, .campo-grupo select {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 15px;
  background-color: #fafafa;
  color: #000000 !important;
}

/* Estilo para la advertencia del teléfono */
.error-texto {
  color: #dc3545;
  font-size: 12px;
  font-weight: 500;
  margin-top: 2px;
}

.modal-botones {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.btn-cancelar {
  flex: 1;
  padding: 12px;
  background-color: #dc3545;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
}

.btn-guardar {
  flex: 1;
  padding: 12px;
  background-color: #28a745;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
}

.btn-guardar:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}
</style>
