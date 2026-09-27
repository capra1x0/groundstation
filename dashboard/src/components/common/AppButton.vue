<template>
  <button
    :type="type"
    class="inline-flex items-center cursor-pointer justify-center gap-2 rounded-md font-medium transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-400 disabled:cursor-not-allowed disabled:opacity-50"
    :class="[variants[variant], sizes[size]]"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    variant?: "primary" | "secondary"
    size?: "sm" | "md"
    type?: "button" | "submit"
  }>(),
  { variant: "primary", size: "md", type: "button" },
)

const variants = {
  primary: "bg-teal-400 text-zinc-950 enabled:hover:bg-teal-300",
  secondary: "border border-zinc-700 text-zinc-200 enabled:hover:bg-zinc-800",
}

const sizes = {
  sm: "px-3 py-1.5 text-sm",
  md: "px-4 py-2 text-sm",
}
</script>