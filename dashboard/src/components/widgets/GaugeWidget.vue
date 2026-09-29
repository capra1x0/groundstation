<template>
  <WidgetCard :widget="widget">
    <div class="flex flex-col items-center">
      <div
        class="relative w-full max-w-56"
        role="meter"
        :aria-valuemin="min"
        :aria-valuemax="max"
        :aria-valuenow="value ?? undefined"
        :aria-label="widget.name"
      >
        <svg viewBox="0 0 200 150" class="w-full" aria-hidden="true">
          <path :d="arc" class="fill-none stroke-zinc-800" stroke-width="14" stroke-linecap="round" />
          <path
            v-if="percent > 0"
            :d="arc"
            :stroke-dasharray="dash"
            class="fill-none stroke-teal-400"
            stroke-width="14"
            stroke-linecap="round"
          />
        </svg>

        <div class="absolute inset-x-0 top-2/3 flex -translate-y-1/2 items-baseline justify-center gap-1">
          <span class="text-3xl font-semibold tabular-nums text-zinc-100">{{ displayValue }}</span>
          <span class="text-sm text-zinc-400">{{ widget.topic.unit }}</span>
        </div>
      </div>

      <div class="flex w-full max-w-56 justify-between px-[9%]">
        <EditableNumber v-model="min" label="Edit minimum" />
        <EditableNumber v-model="max" label="Edit maximum" />
      </div>
    </div>
  </WidgetCard>
</template>

<script setup lang="ts">
import { computed } from "vue"

import EditableNumber from "@/components/common/EditableNumber.vue"
import WidgetCard from "@/components/widgets/WidgetCard.vue"
import { useDashboard } from "@/composables/useDashboard"
import { useReading } from "@/composables/useTelemetry"
import type { DashboardWidget } from "@/types/dashboardWidget"

const props = defineProps<{
  widget: DashboardWidget
}>()

const arc = "M 30.72 140 A 80 80 0 1 1 169.28 140"
const arcLength = 80 * ((4 * Math.PI) / 3)

const { updateSettings } = useDashboard()

const reading = useReading(() => props.widget.topic.topic)

const min = computed({
  get: () => props.widget.settings.min ?? 0,
  set: (newMin: number) => {
    if (newMin < max.value) {
      updateSettings(props.widget.id, { min: newMin })
    }
  },
})

const max = computed({
  get: () => props.widget.settings.max ?? 100,
  set: (newMax: number) => {
    if (newMax > min.value) {
      updateSettings(props.widget.id, { max: newMax })
    }
  },
})

const value = computed(() => (typeof reading.value?.value === "number" ? reading.value.value : null))

const percent = computed(() => {
  if (value.value === null) {
    return 0
  }
  const fraction = (value.value - min.value) / (max.value - min.value)
  return Math.min(Math.max(fraction, 0), 1) * 100
})

const dash = computed(() => `${(percent.value / 100) * arcLength} ${arcLength}`)

const displayValue = computed(() => {
  if (value.value === null) {
    return "–"
  }
  return value.value.toLocaleString(undefined, { minimumFractionDigits: 1, maximumFractionDigits: 1 })
})
</script>