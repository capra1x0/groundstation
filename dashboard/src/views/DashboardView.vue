<template>
  <div
    v-if="widgets.length === 0"
    class="flex h-64 flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-zinc-800 text-center"
  >
    <p class="text-sm text-zinc-300">No widgets yet</p>
    <p class="text-sm text-zinc-500">Add a widget to start watching a topic.</p>
  </div>

  <div v-else class="grid grid-cols-4 gap-4">
    <component
      :is="widgetComponents[entry.widget]"
      v-for="entry in widgets"
      :key="entry.id"
      :widget="entry"
    />
  </div>
</template>

<script setup lang="ts">
import { widgetComponents } from "@/components/widgets"
import { useDashboard } from "@/composables/useDashboard"

const { widgets } = useDashboard()
</script>