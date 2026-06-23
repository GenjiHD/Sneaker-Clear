<template>
  <div class="tarjeta-formulario">
    <h2>Nueva Orden</h2>
    
    <form @submit.prevent="enviarFormulario">
      <div class="grupo-input">
        <input 
          type="text" 
          v-model="nuevaOrden.nombre" 
          @input="validarNombre"
          placeholder="Cliente (Solo letras)" 
          required
        />
      </div>

      <div class="grupo-input">
        <input 
          type="text"
          v-model="nuevaOrden.telefono"
          @input="validarTelefono"
          placeholder="Teléfono (10 dígitos)"
          maxlength="10"
          required
        />
        <small v-if="nuevaOrden.telefono.length > 0 && nuevaOrden.telefono.length < 10" class="error-texto">
          Faltan {{ 10 - nuevaOrden.telefono.length }} dígitos.
        </small>
      </div>

      <div class="grupo-input">
        <input 
          type="text" 
          v-model="nuevaOrden.modelo_marca" 
          placeholder="Marca / Modelo" 
          required
        />
      </div>

      <div class="grupo-input">
        <select v-model="nuevaOrden.tipo_servicio" required>
          <option value="Lavado Básico">Lavado Básico</option>
          <option value="Lavado Premium">Lavado Premium</option>
          <option value="Blanqueado">Blanqueado</option>
        </select>
      </div>

      <div class="grupo-input">
        <input 
          type="number" 
          v-model.number="nuevaOrden.precio" 
          placeholder="Precio" 
          min="0"
          step="1"
          required
        />
      </div>

      <button type="submit" class="btn-guardar" :disabled="nuevaOrden.telefono.length !== 10">
        Guardar Orden
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 1. 🚨 DEFINIMOS EL EMIT PARA CONFIGURAR EL EVENTO PERSONALIZADO
const emit = defineEmits(['orden-creada'])

// Estado del formulario
const nuevaOrden = ref({
  nombre: '',
  telefono: '',
  modelo_marca: '',
  tipo_servicio: 'Lavado Básico', 
  precio: null
})

// Filtro para el nombre
const validarNombre = () => {
  nuevaOrden.value.nombre = nuevaOrden.value.nombre.replace(/[^a-zA-ZáéíóúÁÉÍÓÚñÑüÜ\s]/g, '')
}

// Filtro para el teléfono
const validarTelefono = () => {
  nuevaOrden.value.telefono = nuevaOrden.value.telefono.replace(/\D/g, '')
}

const enviarFormulario = async () => {
  if (nuevaOrden.value.telefono.length !== 10) {
    alert('El teléfono debe tener exactamente 10 dígitos.');
    return;
  }
  if (nuevaOrden.value.precio === null || nuevaOrden.value.precio < 0) {
    alert('El precio no puede ser menor a 0.');
    return;
  }

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
      
      // 2. 🚀 LA ESTOCADA: Le gritamos a la pantalla principal que hay información nueva
      emit('orden-creada');
      
      limpiarFormulario(); 
    } else {
      const errorServidor = await respuesta.text();
      alert(`Error del servidor (500): ${errorServidor || 'No se pudo registrar la orden.'}`);
    }
  } catch (error) {
    console.error('Error de conexion con la API:', error);
    alert('No se pudo conectar con el servidor. Verifica que Go esté encendido.');
  }
}

const limpiarFormulario = () => {
  nuevaOrden.value = {
    nombre: '',
    telefono: '',
    modelo_marca: '',
    tipo_servicio: 'Lavado Básico', 
    precio: null
  }
}
</script>

<style scoped>
.tarjeta-formulario {
  background: #ffffff;
  color: #000000;
  padding: 30px;
  border-radius: 12px;
  max-width: 600px;
  margin: 0 auto 20px auto;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

h2 {
  color: #000000;
  font-size: 24px;
  font-weight: bold;
  margin-top: 0;
  margin-bottom: 25px;
}

.grupo-input {
  margin-bottom: 15px;
}

input, select {
  width: 100%;
  padding: 12px 15px;
  font-size: 16px;
  border: 1px solid #dcdcdc;
  border-radius: 8px;
  box-sizing: border-box;
  background-color: #ffffff;
  
  /* 🍏 CORRECCIÓN: Asegura letras visibles sin importar el tema de NixOS/Hyprland */
  color: #000000 !important;
}

input::placeholder {
  color: #8e8e8e;
}

select {
  background-color: #e9e9e9;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='10' height='10' viewBox='0 0 10 10'><path d='M0 3l5 5 5-5z' fill='%23333'/></svg>");
  background-repeat: no-repeat;
  background-position: right 15px center;
  padding-right: 30px;
}

input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  opacity: 1;
}

/* Letrero de aviso para dígitos faltantes */
.error-texto {
  color: #dc3545;
  font-size: 12px;
  font-weight: 500;
  display: block;
  margin-top: 4px;
}

.btn-guardar {
  width: 100%;
  padding: 14px;
  background-color: #0d47a1; /* Cambiado a un azul primario para mejor contraste */
  color: #ffffff;
  font-weight: bold;
  font-size: 16px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-guardar:hover {
  background-color: #0a357a;
}

.btn-guardar:disabled {
  background-color: #6c757d;
  color: #ffffff;
  cursor: not-allowed;
}
</style>
