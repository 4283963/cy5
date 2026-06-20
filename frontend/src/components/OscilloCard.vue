<script setup lang="ts">
defineProps<{
  title?: string
  subtitle?: string
  class?: string
}>()
</script>

<template>
  <div
    :class="[
      'bg-osc-card border border-osc-border rounded-lg shadow-card relative overflow-hidden',
      $props.class,
    ]"
  >
    <div
      v-if="title"
      class="border-b border-osc-border px-4 py-2.5 flex items-center justify-between"
    >
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 rounded-full bg-osc-green"></div>
        <h3 class="text-osc-bright font-display font-medium text-sm tracking-wide">
          {{ title }}
        </h3>
        <span v-if="subtitle" class="text-osc-muted text-xs">{{ subtitle }}</span>
      </div>
    </div>
    <div class="p-4">
      <slot />
    </div>
  </div>
</template>
