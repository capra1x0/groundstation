<template>
  <header class="flex items-center gap-4 border-b border-zinc-800 px-6 py-3">
    <RouterLink to="/" class="flex items-center gap-2 font-semibold">
      <span class="h-2.5 w-2.5 rotate-45 bg-teal-400"></span>
      Groundstation
    </RouterLink>

    <span class="flex items-center gap-2 rounded-md border border-zinc-800 px-2.5 py-1 font-mono text-xs text-zinc-400">
      <span class="h-1.5 w-1.5 rounded-full" :class="connected ? 'bg-teal-400' : 'bg-red-500'"></span>
      {{ connected ? "Live" : "Backend offline" }}
    </span>

    <span class="ml-auto font-mono text-xs text-zinc-500">
      {{ topics.length }} topics · 0 widgets
    </span>

    <Button label="+ Add widget" size="small" />
  </header>
</template>

<script setup lang="ts">
import Button from "primevue/button"

import { useTelemetry } from "@/composables/useTelemetry"

const { connected, topics } = useTelemetry()
</script>