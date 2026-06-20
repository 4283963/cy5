<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'

const { applyTheme } = useTheme()

onMounted(() => {
  applyTheme('dark')
})
</script>

<template>
  <div class="min-h-screen bg-osc-bg text-osc-text">
    <router-view />
  </div>
</template>
