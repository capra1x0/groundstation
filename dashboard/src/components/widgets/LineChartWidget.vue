<template>
  <WidgetCard :widget="widget" class="col-span-2">
    <div class="flex flex-col gap-3">
      <div class="flex items-baseline gap-2">
        <span class="text-3xl font-semibold tabular-nums text-zinc-100">{{ displayValue }}</span>
        <span class="text-sm text-zinc-400">{{ widget.topic.unit }}</span>
      </div>

      <div
        class="relative h-40 cursor-ew-resize touch-none select-none border-l border-zinc-800"
        title="Drag or scroll to change the time range, double-click to reset"
        @pointerdown="startDrag"
        @pointermove="drag"
        @pointerup="endDrag"
        @pointercancel="endDrag"
        @wheel.prevent="zoom"
        @dblclick="resetWindow"
      >
        <svg :viewBox="`0 0 ${WIDTH} ${HEIGHT}`" preserveAspectRatio="none" class="h-full w-full" aria-hidden="true">
          <line
            v-for="y in gridLines"
            :key="y"
            x1="0"
            :x2="WIDTH"
            :y1="y"
            :y2="y"
            class="stroke-zinc-800"
            stroke-dasharray="4 4"
            vector-effect="non-scaling-stroke"
          />
          <template v-if="points.length > 1">
            <path :d="areaPath" class="fill-teal-400/10" />
            <path
              :d="linePath"
              class="fill-none stroke-teal-400"
              stroke-width="2"
              stroke-linejoin="round"
              vector-effect="non-scaling-stroke"
            />
          </template>
        </svg>

        <template v-if="points.length > 1">
          <span class="absolute left-2 top-1 font-mono text-xs text-zinc-500">{{ formatValue(range.max) }}</span>
          <span class="absolute bottom-1 left-2 font-mono text-xs text-zinc-500">{{ formatValue(range.min) }}</span>
        </template>
        <p v-else class="absolute inset-0 flex items-center justify-center text-sm text-zinc-500">Waiting for data</p>
      </div>

      <div class="flex justify-between font-mono text-xs text-zinc-500">
        <span>-{{ formatDuration(windowSeconds) }}</span>
        <span>-{{ formatDuration(windowSeconds / 2) }}</span>
        <span>now</span>
      </div>
    </div>
  </WidgetCard>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue"

import WidgetCard from "@/components/widgets/WidgetCard.vue"
import { useDashboard } from "@/composables/useDashboard"
import { getHistory, useReading } from "@/composables/useTelemetry"
import type { DashboardWidget } from "@/types/dashboardWidget"

interface Point {
  ts: number
  value: number
}

const WIDTH = 600
const HEIGHT = 160
const PADDING = 8

const MIN_WINDOW = 10
const MAX_WINDOW = 300
const DEFAULT_WINDOW = 60

const gridLines = [HEIGHT / 3, (HEIGHT / 3) * 2]

const props = defineProps<{
  widget: DashboardWidget
}>()

const { updateSettings } = useDashboard()
const reading = useReading(() => props.widget.topic.topic)

const now = ref(Date.now())
let timer: number | undefined

onMounted(() => {
  timer = window.setInterval(() => {
    now.value = Date.now()
  }, 200)
})

onUnmounted(() => {
  window.clearInterval(timer)
})

const windowSeconds = computed(() => props.widget.settings.windowSeconds ?? DEFAULT_WINDOW)

function setWindow(seconds: number) {
  const clamped = Math.min(Math.max(seconds, MIN_WINDOW), MAX_WINDOW)
  updateSettings(props.widget.id, { windowSeconds: clamped })
}

function resetWindow() {
  setWindow(DEFAULT_WINDOW)
}

const dragging = ref(false)
let dragStartX = 0
let dragStartWindow = 0

function startDrag(event: PointerEvent) {
  dragging.value = true
  dragStartX = event.clientX
  dragStartWindow = windowSeconds.value

  const element = event.currentTarget as HTMLElement
  element.setPointerCapture(event.pointerId)
}

function drag(event: PointerEvent) {
  if (!dragging.value) {
    return
  }

  const distance = event.clientX - dragStartX
  setWindow(dragStartWindow * Math.pow(2, -distance / 200))
}

function endDrag() {
  dragging.value = false
}

function zoom(event: WheelEvent) {
  setWindow(windowSeconds.value * Math.pow(2, event.deltaY / 500))
}

const start = computed(() => now.value - windowSeconds.value * 1000)

const points = computed(() => {
  const history = getHistory(props.widget.topic.topic)

  const firstVisible = history.findIndex((entry) => entry.ts >= start.value)
  if (firstVisible === -1) {
    return []
  }

  const visible: Point[] = []
  for (const entry of history.slice(Math.max(firstVisible - 1, 0))) {
    if (typeof entry.value === "number") {
      visible.push({ ts: entry.ts, value: entry.value })
    }
  }

  return downsample(visible)
})

function downsample(input: Point[]): Point[] {
  if (input.length <= WIDTH * 2) {
    return input
  }

  const columnMs = (windowSeconds.value * 1000) / WIDTH
  const result: Point[] = []
  let column: Point[] = []
  let columnIndex: number | null = null

  for (const point of input) {
    const index = Math.floor((point.ts - start.value) / columnMs)
    if (index !== columnIndex) {
      result.push(...lowestAndHighest(column))
      column = []
      columnIndex = index
    }
    column.push(point)
  }
  result.push(...lowestAndHighest(column))

  return result
}

function lowestAndHighest(column: Point[]): Point[] {
  if (column.length <= 2) {
    return column
  }

  const lowest = column.reduce((low, point) => (point.value < low.value ? point : low))
  const highest = column.reduce((high, point) => (point.value > high.value ? point : high))

  return lowest.ts <= highest.ts ? [lowest, highest] : [highest, lowest]
}

const range = computed(() => {
  let min = Infinity
  let max = -Infinity

  for (const point of points.value) {
    min = Math.min(min, point.value)
    max = Math.max(max, point.value)
  }

  if (min === max) {
    min -= 1
    max += 1
  }

  return { min, max }
})

function toX(ts: number): number {
  return ((ts - start.value) / (windowSeconds.value * 1000)) * WIDTH
}

function toY(value: number): number {
  const fraction = (value - range.value.min) / (range.value.max - range.value.min)
  return HEIGHT - PADDING - fraction * (HEIGHT - PADDING * 2)
}

const linePath = computed(() =>
  points.value
    .map((point, index) => `${index === 0 ? "M" : "L"} ${toX(point.ts).toFixed(1)} ${toY(point.value).toFixed(1)}`)
    .join(" "),
)

const areaPath = computed(() => {
  const first = points.value[0]
  const last = points.value[points.value.length - 1]
  if (!first || !last) {
    return ""
  }
  return `${linePath.value} L ${toX(last.ts).toFixed(1)} ${HEIGHT} L ${toX(first.ts).toFixed(1)} ${HEIGHT} Z`
})

const displayValue = computed(() => {
  const value = reading.value?.value
  return typeof value === "number" ? formatValue(value) : "–"
})

function formatValue(value: number): string {
  return value.toLocaleString(undefined, { maximumFractionDigits: 1 })
}

function formatDuration(seconds: number): string {
  if (seconds < 120) {
    return `${Math.round(seconds)}s`
  }
  return `${(seconds / 60).toLocaleString(undefined, { maximumFractionDigits: 1 })}m`
}
</script>