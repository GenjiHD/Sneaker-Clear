<template>
  <div class="contenedor-principal">
    
    <!-- 🏢 HEADER SUPERIOR CON LOGO DINÁMICO -->
    <header class="header-tpv">
      <div class="logo-titulo">
        <!-- Renderiza la imagen cargada o el emoji por defecto -->
        <div class="contenedor-logo-dinamico">
          <img v-if="urlLogo" :src="urlLogo" alt="Logo TPV" class="logo-subido" />
          <span v-else class="icono-logo">👟</span>
        </div>
        
        <div class="texto-header">
          <h1>Sneaker Cleaning TPV PRO</h1>
          <p>Sistema Profesional de Control de Órdenes</p>
        </div>
      </div>

      <!-- Input de archivo oculto controlado por el botón -->
      <input 
        type="file" 
        ref="inputArchivo" 
        @change="procesarNuevoLogo" 
        accept="image/*" 
        style="display: none;" 
      />
      <button class="btn-logo" @click="seleccionarImagen">Cambiar Logotipo</button>
    </header>

    <!-- 🛠️ EL CUERPO EN DOS COLUMNAS -->
    <main class="grid-tpv">
      <section class="columna-formulario">
        <!-- Escuchamos el evento de orden creada para refrescar la tabla y métricas -->
        <FormularioServicio @orden-creada="refrescarDatosGlobales" />
      </section>

      <section class="columna-contenido">
        <!-- Pasamos la referencia para disparar la actualización desde aquí -->
        <ListaOrdenes ref="componenteLista" />
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import FormularioServicio from './components/FormularioServicio.vue'
import ListaOrdenes from './components/ListaOrdenes.vue'

const inputArchivo = ref(null)
const urlLogo = ref(null)
const componenteLista = ref(null)

// Disparar clic al input oculto
const seleccionarImagen = () => {
  inputArchivo.value.click()
}

// Convertir imagen a Base64 de forma asíncrona y guardarla en LocalStorage
const procesarNuevoLogo = (evento) => {
  const archivo = evento.target.files[0]
  if (!archivo) return

  const lector = new FileReader()
  lector.onload = (e) => {
    const base64String = e.target.result
    urlLogo.value = base64String
    localStorage.setItem('tpv_logo_custom', base64String)
  }
  lector.readAsDataURL(archivo)
}

// Escucha cuando el formulario agrega algo y le avisa a la tabla que se actualice
const refrescarDatosGlobales = () => {
  if (componenteLista.value && componenteLista.value.obtenerOrdenes) {
    componenteLista.value.obtenerOrdenes()
  }
}

onMounted(() => {
  // Cargar logotipo persistido si existe
  const logoGuardado = localStorage.getItem('tpv_logo_custom')
  if (logoGuardado) {
    urlLogo.value = logoGuardado
  }
})
</script>

<style scoped>
.contenedor-principal {
  min-height: 100vh;
  background-color: #f3f4f6;
  padding: 24px;
  box-sizing: border-box;
  font-family: system-ui, -apple-system, sans-serif;
}

.header-tpv {
  background-color: #ffffff;
  padding: 16px 24px;
  border-radius: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.02);
  border: 1px solid #e5e7eb;
}

.logo-titulo {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* Manejo del tamaño del contenedor del logo */
.contenedor-logo-dinamico {
  width: 46px;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-subido {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.icono-logo {
  font-size: 24px;
  background: #f1f5f9;
  padding: 10px;
  border-radius: 12px;
  display: block;
}

.texto-header h1 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  color: #0f172a;
}

.texto-header p {
  margin: 2px 0 0 0;
  font-size: 13px;
  color: #64748b;
}

.btn-logo {
  background-color: #0f172a;
  color: #ffffff;
  border: none;
  padding: 10px 16px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-logo:hover {
  background-color: #1e293b;
}

.grid-tpv {
  display: flex;
  gap: 24px;
  align-items: flex-start;
  max-width: 1400px;
  margin: 0 auto;
}

.columna-formulario {
  flex: 0 0 360px;
}

.columna-contenido {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

@media (max-width: 1024px) {
  .grid-tpv { flex-direction: column; }
  .columna-formulario { width: 100%; flex: none; }
}
</style>