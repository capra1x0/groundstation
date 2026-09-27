<template>
  <div class="min-h-screen bg-zinc-950 text-zinc-100">
    <AppHeader />

    <main class="p-6">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import AppHeader from "@/components/layout/AppHeader.vue"
</script>