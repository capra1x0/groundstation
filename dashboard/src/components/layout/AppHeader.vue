<template>
  <header class="flex items-center gap-4 border-b border-zinc-800 px-6 py-3">
    <RouterLink to="/" class="flex items-center gap-2 font-semibold">
      <span class="h-2.5 w-2.5 rotate-45 bg-teal-400"></span>
      groundstation
    </RouterLink>

    <span class="flex items-center gap-2 rounded-md border border-zinc-800 px-2.5 py-1 font-mono text-xs text-zinc-400">
      <span class="h-1.5 w-1.5 rounded-full" :class="connected ? 'bg-teal-400' : 'bg-red-500'"></span>
      {{ connected ? "Live" : "Backend offline" }}
    </span>

    <span class="ml-auto font-mono text-xs text-zinc-500">
      {{ topics.length }} topics · {{ widgets.length }} widgets
    </span>

    <AppButton size="sm" @click="addWidgetOpen = true">+ Add widget</AppButton>

    <AddWidgetDialog v-model:visible="addWidgetOpen" @add="addWidget" />
  </header>
</template>

<script setup lang="ts">
import { ref } from "vue"

import AppButton from "@/components/common/AppButton.vue"
import AddWidgetDialog from "@/components/dialogs/AddWidgetDialog.vue"
import { useDashboard } from "@/composables/useDashboard"
import { useTelemetry } from "@/composables/useTelemetry"

const { connected, topics } = useTelemetry()
const { widgets, addWidget } = useDashboard()

const addWidgetOpen = ref(false)
</script>