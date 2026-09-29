<template>
  <input
    v-if="editing"
    ref="inputElement"
    v-model="draft"
    type="number"
    step="any"
    class="w-16 rounded border border-teal-400 bg-zinc-950 px-1 text-center font-mono text-xs text-zinc-100 focus:outline-none"
    :aria-label="label"
    @keydown.enter="save"
    @keydown.escape="cancel"
    @blur="save"
  />
  <button
    v-else
    type="button"
    class="rounded px-1 font-mono text-xs text-zinc-500 hover:bg-zinc-800 hover:text-zinc-100"
    :aria-label="label"
    @click="startEditing"
  >
    {{ model }}
  </button>
</template>

<script setup lang="ts">
import { nextTick, ref } from "vue"

defineProps<{
  label: string
}>()

const model = defineModel<number>({ required: true })

const editing = ref(false)
const draft = ref<string | number>("")
const inputElement = ref<HTMLInputElement | null>(null)

async function startEditing() {
  draft.value = model.value
  editing.value = true

  await nextTick()
  inputElement.value?.select()
}

function save() {
  if (!editing.value) {
    return
  }
  editing.value = false

  const text = String(draft.value).trim()
  const value = Number(text)
  if (text !== "" && Number.isFinite(value)) {
    model.value = value
  }
}

function cancel() {
  editing.value = false
}
</script>