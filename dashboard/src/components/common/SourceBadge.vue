<template>
  <span class="rounded border px-1.5 py-px font-mono text-xs" :class="colorClasses">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from "vue"

const props = defineProps<{
  label: string
}>()

const colors: Record<string, string> = {
  ESP32: "border-teal-400/60 text-teal-300",
  Simulation: "border-violet-400/60 text-violet-300",
}

const colorClasses = computed(() => colors[props.label] ?? "border-amber-400/60 text-amber-300")
</script>