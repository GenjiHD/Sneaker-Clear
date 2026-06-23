<template>
  <div class="contenedor-principal">
    <header>
      <h1>👟 Sneaker Cleaning - Panel de Control</h1>
    </header>

    <main>
      <FormularioServicio />

      <ListaOrdenes />

    </main>
  </div>
</template>

<script setup>
  import FormularioServicio from './components/FormularioServicio.vue'
  import ListaOrdenes from './components/ListaOrdenes.vue'
</script>


