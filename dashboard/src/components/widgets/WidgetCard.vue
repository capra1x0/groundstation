<template>
  <article class="flex flex-col rounded-xl border border-zinc-800 bg-zinc-900 p-4">
    <header class="flex items-start justify-between gap-3">
      <div class="flex min-w-0 flex-col gap-0.5">
        <h2 class="truncate font-semibold text-zinc-100">{{ widget.name }}</h2>
        <p v-if="widget.description" class="text-sm text-zinc-400">{{ widget.description }}</p>
      </div>

      <button
        type="button"
        class="text-lg leading-none text-zinc-500 hover:text-zinc-100"
        aria-label="Remove widget"
        @click="removeWidget(widget.id)"
      >
        ×
      </button>
    </header>

    <div class="flex-1 py-4">
      <slot />
    </div>

    <footer class="flex items-center gap-2 border-t border-zinc-800 pt-3 text-xs">
      <SourceBadge :label="widget.topic.sourceLabel" />
      <span class="min-w-0 flex-1 truncate font-mono text-zinc-500">{{ widget.topic.topic }}</span>
      <span class="font-mono text-zinc-500">{{ lastUpdate }}</span>
    </footer>
  </article>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { format } from "date-fns"

import SourceBadge from "@/components/common/SourceBadge.vue"
import { useDashboard } from "@/composables/useDashboard"
import { useReading } from "@/composables/useTelemetry"
import type { DashboardWidget } from "@/types/dashboardWidget"

const props = defineProps<{
  widget: DashboardWidget
}>()

const { removeWidget } = useDashboard()

const reading = useReading(() => props.widget.topic.topic)

const lastUpdate = computed(() => (reading.value ? format(reading.value.ts, "HH:mm:ss") : "–"))
</script>