<template>
  <AppDialog v-model:visible="visible" aria-labelledby="add-widget-title">
    <div class="flex max-h-[85vh] flex-col">
      <header class="flex flex-col gap-4 border-b border-zinc-800 px-6 py-5">
        <div class="flex items-center justify-between">
          <h2 id="add-widget-title" class="text-lg font-semibold text-zinc-100">Add widget</h2>
          <button
            type="button"
            class="text-xl leading-none text-zinc-400 hover:text-zinc-100 cursor-pointer"
            aria-label="Close"
            @click="visible = false"
          >
            ×
          </button>
        </div>

        <ol class="flex items-center gap-3 text-sm">
          <li class="flex items-center gap-2" :class="step === 1 ? 'text-teal-400' : 'text-zinc-400'">
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full border text-xs"
              :class="step === 1 ? 'border-teal-400' : 'border-zinc-500'"
            >
              1
            </span>
            Topic
          </li>
          <li class="h-px w-12 bg-zinc-700" aria-hidden="true"></li>
          <li class="flex items-center gap-2" :class="step === 2 ? 'text-teal-400' : 'text-zinc-400'">
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full border text-xs"
              :class="step === 2 ? 'border-teal-400' : 'border-zinc-500'"
            >
              2
            </span>
            Display
          </li>
        </ol>
      </header>

      <!-- Step 1 -->
      <div v-if="step === 1" class="min-h-0 flex-1 overflow-y-auto p-3">
        <p v-if="loading" class="px-3 py-8 text-center text-sm text-zinc-500">Loading topics…</p>

        <div v-else-if="error" class="flex flex-col items-center gap-3 px-3 py-8 text-center">
          <p class="text-sm text-zinc-400">Could not load topics. Check that the telemetry backend is running.</p>
          <AppButton variant="secondary" size="sm" @click="loadTopics">Try again</AppButton>
        </div>

        <p v-else-if="topics.length === 0" class="px-3 py-8 text-center text-sm text-zinc-400">
          No topics defined yet. Add them to telemetry/topics.json.
        </p>

        <div v-else role="radiogroup" aria-label="Topic" class="flex flex-col gap-1">
          <button
            v-for="topic in topics"
            :key="topic.topic"
            type="button"
            role="radio"
            :aria-checked="selectedId === topic.topic"
            class="flex w-full items-start gap-4 rounded-lg border px-4 py-3 text-left transition-colors"
            :class="selectedId === topic.topic ? 'border-teal-400/70 bg-teal-400/10' : 'border-transparent hover:bg-zinc-800/50'"
            @click="selectedId = topic.topic"
          >
            <span
              class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full border"
              :class="selectedId === topic.topic ? 'border-teal-400' : 'border-zinc-600'"
            >
              <span v-if="selectedId === topic.topic" class="h-2.5 w-2.5 rounded-full bg-teal-400"></span>
            </span>

            <span class="flex min-w-0 flex-1 flex-col gap-1">
              <span class="flex items-center gap-2">
                <span class="font-semibold text-zinc-100">{{ topic.name }}</span>
                <SourceBadge :label="topic.sourceLabel" />
              </span>
              <span class="truncate font-mono text-xs text-zinc-500">{{ topic.topic }}</span>
              <span class="text-sm text-zinc-400">{{ topic.description }}</span>
            </span>
          </button>
        </div>
      </div>

      <!-- Step 2 -->
      <div v-else class="flex min-h-0 flex-1 flex-col gap-6 overflow-y-auto px-6 py-5">
        <div class="flex items-center justify-between gap-4 rounded-lg bg-zinc-800/50 px-4 py-3">
          <div class="flex min-w-0 flex-col gap-1">
            <span class="font-semibold text-zinc-100">{{ selectedTopic?.name }}</span>
            <span class="truncate font-mono text-xs text-zinc-500">{{ selectedTopic?.topic }}</span>
          </div>
          <button type="button" class="cursor-pointer shrink-0 text-sm text-teal-400 hover:text-teal-300" @click="step = 1">
            Change
          </button>
        </div>

        <section class="flex flex-col gap-3">
          <h3 class="text-sm">
            <span class="font-semibold text-zinc-100">Display component</span>
            <span class="ml-2 text-xs text-zinc-500">Suited to {{ selectedTopic?.valueType }} values</span>
          </h3>

          <div
            role="radiogroup"
            aria-label="Display component"
            class="flex snap-x gap-3 overflow-x-auto pb-2"
          >
            <button
              v-for="option in widgetOptions"
              :key="option.id"
              type="button"
              role="radio"
              :disabled="!isAvailable(option)"
              :aria-checked="selectedDisplay === option.id"
              class="flex aspect-11/10 w-44 shrink-0 snap-start flex-col gap-2 rounded-lg border p-3 text-left transition-colors disabled:cursor-not-allowed disabled:opacity-40"
              :class="selectedDisplay === option.id ? 'border-teal-400/70 bg-teal-400/10' : 'border-zinc-700 enabled:hover:bg-zinc-800/50'"
              @click="selectedDisplay = option.id"
            >
              <span class="flex flex-1 items-center justify-center">
                <img :src="option.preview" alt="" class="max-h-full" />
              </span>
              <span class="flex flex-col gap-0.5">
                <span class="text-sm font-semibold text-zinc-100">{{ option.name }}</span>
                <span class="line-clamp-2 text-xs text-zinc-400">{{ option.description }}</span>
              </span>
            </button>
          </div>
        </section>

        <div class="grid gap-4 sm:grid-cols-2">
          <div class="flex flex-col gap-2">
            <label for="widget-name" class="text-sm font-semibold text-zinc-100">Name</label>
            <input
              id="widget-name"
              v-model="widgetName"
              type="text"
              class="w-full rounded-md border border-zinc-700 bg-zinc-950 px-3 py-2 text-sm text-zinc-100 focus:border-teal-400 focus:outline-none"
            />
          </div>

          <div class="flex flex-col gap-2">
            <label for="widget-description" class="text-sm">
              <span class="font-semibold text-zinc-100">Description</span>
              <span class="ml-1 text-zinc-500">optional</span>
            </label>
            <textarea
              id="widget-description"
              v-model="widgetDescription"
              rows="2"
              placeholder="Shown under the title"
              class="w-full resize-y rounded-md border border-zinc-700 bg-zinc-950 px-3 py-2 text-sm text-zinc-100 placeholder:text-zinc-500 focus:border-teal-400 focus:outline-none"
            ></textarea>
          </div>
        </div>
      </div>

      <footer class="flex items-center gap-3 border-t border-zinc-800 px-6 py-4">
        <template v-if="step === 1">
          <span class="min-w-0 flex-1 truncate font-mono text-xs text-zinc-500">{{ selectedId }}</span>
          <AppButton variant="secondary" @click="visible = false">Cancel</AppButton>
          <AppButton :disabled="!selectedTopic" @click="goToDisplay">Next</AppButton>
        </template>

        <template v-else>
          <span class="flex-1"></span>
          <AppButton variant="secondary" @click="step = 1">Back</AppButton>
          <AppButton :disabled="!canAdd" @click="addWidget">Add to dashboard</AppButton>
        </template>
      </footer>
    </div>
  </AppDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from "vue"

import { getTopics } from "@/api/topics"
import AppButton from "@/components/common/AppButton.vue"
import AppDialog from "@/components/common/AppDialog.vue"
import SourceBadge from "@/components/common/SourceBadge.vue"
import type { Topic } from "@/types/topic"
import { widgetOptions, type WidgetOption } from "@/config/widgets"

const visible = defineModel<boolean>("visible", { required: true })

const emit = defineEmits<{
  add: [widget: { topic: string, display: string, name: string, description: string }]
}>()

const step = ref<1 | 2>(1)

// Step 1
const topics = ref<Topic[]>([])
const selectedId = ref<string | null>(null)
const loading = ref(false)
const error = ref(false)

const selectedTopic = computed(() => topics.value.find((topic) => topic.topic === selectedId.value))

// Step 2
const selectedDisplay = ref<string | null>(null)
const widgetName = ref("")
const widgetDescription = ref("")

const canAdd = computed(() => selectedDisplay.value !== null && widgetName.value.trim() !== "")

async function loadTopics() {
  loading.value = true
  error.value = false

  try {
    topics.value = await getTopics()
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

function isAvailable(option: WidgetOption): boolean {
  if(!selectedTopic.value) {
    return false
  }
  return option.valueTypes.includes(selectedTopic.value.valueType)
}

watch(visible, (open) => {
  if (open) {
    step.value = 1
    selectedId.value = null
    loadTopics()
  }
})

function goToDisplay() {
  if (!selectedTopic.value) {
    return
  }

  selectedDisplay.value = widgetOptions.find(isAvailable)?.id ?? null
  widgetName.value = selectedTopic.value.name
  widgetDescription.value = ""
  step.value = 2
}

function addWidget() {
  if (!selectedTopic.value || !selectedDisplay.value) {
    return
  }

  emit("add", {
    topic: selectedTopic.value.topic,
    display: selectedDisplay.value,
    name: widgetName.value.trim(),
    description: widgetDescription.value.trim(),
  })

  visible.value = false
}
</script>