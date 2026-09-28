<template>
  <WidgetCard :widget="widget">
    <div class="flex items-baseline gap-2">
      <span class="text-5xl font-semibold tabular-nums text-zinc-100">{{ displayValue }}</span>
      <span class="text-lg text-zinc-400">{{ widget.topic.unit }}</span>
    </div>
  </WidgetCard>
</template>

<script setup lang="ts">
import { computed } from "vue"

import WidgetCard from "@/components/widgets/WidgetCard.vue"
import { useReading } from "@/composables/useTelemetry"
import type { DashboardWidget } from "@/types/dashboardWidget"

const props = defineProps<{
  widget: DashboardWidget
}>()

const reading = useReading(() => props.widget.topic.topic)

const displayValue = computed(() => {
  const value = reading.value?.value
  if (typeof value !== "number") {
    return "–"
  }
  return value.toLocaleString(undefined, { maximumFractionDigits: 1, minimumFractionDigits: 1 })
})
</script>