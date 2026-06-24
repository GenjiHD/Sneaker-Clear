<template>
  <div class="tarjeta-formulario-pro">
    <h2 class="titulo-formulario">Nueva Orden</h2>
    
    <form @submit.prevent="enviarFormulario" class="cuerpo-formulario">
      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Cliente</label>
        <input 
          type="text" 
          v-model="nuevaOrden.nombre" 
          @input="validarNombre"
          placeholder="Nombre del cliente" 
          class="control-input-pro"
          required
        />
      </div>

      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Teléfono</label>
        <input 
          type="text"
          v-model="nuevaOrden.telefono"
          @input="validarTelefono"
          placeholder="6671234567"
          maxlength="10"
          class="control-input-pro"
          required
        />
        <small v-if="nuevaOrden.telefono.length > 0 && nuevaOrden.telefono.length < 10" class="error-texto-pro">
          Faltan {{ 10 - nuevaOrden.telefono.length }} dígitos.
        </small>
      </div>

      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Marca / Modelo</label>
        <input 
          type="text" 
          v-model="nuevaOrden.modelo_marca" 
          placeholder="Nike Air Force 1" 
          class="control-input-pro"
          required
        />
      </div>

      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Servicio</label>
        <div class="contenedor-select-pro">
          <select v-model="nuevaOrden.tipo_servicio" class="control-select-pro" required>
            <option value="Lavado Básico">Lavado Básico</option>
            <option value="Lavado Premium">Lavado Premium</option>
            <option value="Blanqueado">Blanqueado</option>
          </select>
        </div>
      </div>

      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Precio</label>
        <input 
          type="number" 
          v-model.number="nuevaOrden.precio" 
          placeholder="0" 
          min="0"
          step="1"
          class="control-input-pro"
          required
        />
      </div>

      <div class="grupo-input-pro">
        <label class="etiqueta-pro">Notas</label>
        <textarea 
          v-model="nuevaOrden.observaciones" 
          placeholder="Detalles del servicio"
          rows="3"
          class="control-textarea-pro"
        ></textarea>
      </div>

      <div class="fila-botones-pro">
        <button type="submit" class="btn-guardar-pro" :disabled="nuevaOrden.telefono.length !== 10">
          Guardar Orden
        </button>
        <button type="button" class="btn-limpiar-pro" @click="limpiarFormulario">
          Limpiar
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['orden-creada'])

const nuevaOrden = ref({
  nombre: '',
  telefono: '',
  modelo_marca: '',
  tipo_servicio: 'Lavado Básico', 
  precio: null,
  observaciones: '' 
})

const validarNombre = () => {
  nuevaOrden.value.nombre = nuevaOrden.value.nombre.replace(/[^a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]/g, '')
}

const validarTelefono = () => {
  nuevaOrden.value.telefono = nuevaOrden.value.telefono.replace(/\D/g, '')
}

const enviarFormulario = async () => {
  if (nuevaOrden.value.telefono.length !== 10) return;
  if (nuevaOrden.value.precio === null || nuevaOrden.value.precio < 0) return;

  nuevaOrden.value.precio = parseFloat(nuevaOrden.value.precio);
  const urlAPI = import.meta.env.VITE_API_URL;

  try {
    const respuesta = await fetch(`${urlAPI}/servicios/crear`, {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(nuevaOrden.value)
    });

    if (respuesta.ok) {
      alert(`¡Orden de ${nuevaOrden.value.nombre} guardada con éxito!`);
      emit('orden-creada');
      limpiarFormulario(); 
    } else {
      const errorServidor = await respuesta.text();
      alert(`Error del servidor (500): ${errorServidor || 'No se pudo registrar la orden.'}`);
    }
  } catch (error) {
    console.error('Error de conexion con la API:', error);
    alert('No se pudo conectar con el servidor.');
  }
}

const limpiarFormulario = () => {
  nuevaOrden.value = {
    nombre: '',
    telefono: '',
    modelo_marca: '',
    tipo_servicio: 'Lavado Básico', 
    precio: null,
    observaciones: '' 
  }
}
</script>

<style scoped>
/* Contenedor principal de la tarjeta izquierda */
.tarjeta-formulario-pro {
  background: #ffffff;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.03), 0 1px 3px rgba(0, 0, 0, 0.01);
  width: 100%;
  max-width: 400px; /* Ajustado para que mantenga la proporción esbelta de la captura */
  box-sizing: border-box;
  border: 1px solid #f3f4f6;
}

/* Título superior h2 */
.titulo-formulario {
  color: #0f172a;
  font-size: 20px;
  font-weight: 700;
  margin-top: 0;
  margin-bottom: 24px;
  letter-spacing: -0.02em;
}

.cuerpo-formulario {
  display: flex;
  flex-direction: column;
  gap: 14px; /* Espaciado uniforme entre bloques */
}

.grupo-input-pro {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

/* Estilo exacto de las etiquetas del prototipo */
.etiqueta-pro {
  font-size: 13px;
  font-weight: 600;
  color: #1e293b;
  letter-spacing: -0.01em;
}

/* Estilización global de los campos de texto y número */
.control-input-pro,
.control-textarea-pro,
.control-select-pro {
  width: 100%;
  padding: 11px 14px;
  font-size: 14px;
  font-weight: 400;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  box-sizing: border-box;
  background-color: #f8fafc; /* Color de fondo sutil característico */
  color: #0f172a !important;
  transition: all 0.2s ease;
}

/* Efecto focus limpio sin outlines toscos */
.control-input-pro:focus,
.control-textarea-pro:focus,
.control-select-pro:focus {
  outline: none;
  border-color: #cbd5e1;
  background-color: #ffffff;
}

/* Placeholders finos */
.control-input-pro::placeholder,
.control-textarea-pro::placeholder {
  color: #94a3b8;
  font-weight: 400;
}

/* Área de notas fija */
.control-textarea-pro {
  resize: none;
  font-family: inherit;
  line-height: 1.5;
}

/* Contenedor del select para la flecha personalizada */
.contenedor-select-pro {
  position: relative;
  width: 100%;
}

.control-select-pro {
  cursor: pointer;
  appearance: none;
  padding-right: 36px;
}

/* Icono SVG de flecha hacia abajo idéntico al del prototipo */
.contenedor-select-pro::after {
  content: "";
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  width: 10px;
  height: 6px;
  background-repeat: no-repeat;
  background-size: contain;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='10' height='6' viewBox='0 0 10 6' fill='none'><path d='M1 1L5 5L9 1' stroke='%23334155' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/></svg>");
  pointer-events: none;
}

/* Texto de error para dígitos del teléfono */
.error-texto-pro {
  color: #ef4444;
  font-size: 11px;
  font-weight: 500;
  margin-top: 1px;
}

/* Contenedor flex para botones alineados */
.fila-botones-pro {
  display: flex;
  gap: 12px;
  margin-top: 10px;
}

/* Botón principal Guardar Orden (Grande y oscuro) */
.btn-guardar-pro {
  flex: 1.8;
  padding: 12px 16px;
  background-color: #0f172a;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.btn-guardar-pro:hover {
  background-color: #1e293b;
}

.btn-guardar-pro:disabled {
  background-color: #cbd5e1;
  color: #94a3b8;
  cursor: not-allowed;
}

/* Botón secundario Limpiar (Pequeño y gris intermedio) */
.btn-limpiar-pro {
  flex: 1;
  padding: 12px 16px;
  background-color: #334155;
  color: #ffffff;
  font-weight: 600;
  font-size: 14px;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.btn-limpiar-pro:hover {
  background-color: #475569;
}
</style>
