<template>
  <dialog
    ref="dialogElement"
    class="m-auto w-[calc(100%-2rem)] max-w-3xl rounded-xl border border-zinc-800 bg-zinc-900 p-0 text-zinc-100 backdrop:bg-black/60"
    @close="visible = false"
    @click="closeOnBackdrop"
  >
    <slot />
  </dialog>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue"

const visible = defineModel<boolean>("visible", { required: true })

const dialogElement = ref<HTMLDialogElement | null>(null)

function sync(open: boolean) {
  if (open) {
    dialogElement.value?.showModal()
  } else {
    dialogElement.value?.close()
  }
}

watch(visible, sync)
onMounted(() => sync(visible.value))

function closeOnBackdrop(event: MouseEvent) {
  if (event.target === dialogElement.value) {
    visible.value = false
  }
}
</script>