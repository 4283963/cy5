<script setup lang="ts">
import { MapPin } from 'lucide-vue-next'
import type { Peak } from '@/types'

defineProps<{
  peaks: Peak[]
  highlightIndex?: number
}>()

const emit = defineEmits<{
  (e: 'highlight', index: number): void
}>()
</script>

<template>
  <div class="space-y-2 max-h-[240px] overflow-y-auto osc-scrollbar pr-1">
    <div v-if="peaks.length === 0" class="text-center py-8 text-osc-muted text-sm">
      暂无检测到的缺陷峰
    </div>
    <div
      v-for="(peak, idx) in peaks"
      :key="idx"
      @click="emit('highlight', idx)"
      :class="[
        'group px-3 py-2 rounded-md cursor-pointer transition-all duration-150 border',
        highlightIndex === idx
          ? 'bg-osc-red/10 border-osc-red/40'
          : 'bg-osc-bg/50 border-transparent hover:bg-osc-bg hover:border-osc-border',
      ]"
    >
      <div class="flex items-center gap-3">
        <div
          :class="[
            'w-6 h-6 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-bold',
            highlightIndex === idx
              ? 'bg-osc-red text-white shadow-red-glow'
              : 'bg-osc-red/20 text-osc-red group-hover:bg-osc-red/30',
          ]"
        >
          {{ idx + 1 }}
        </div>
        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between">
            <span class="text-osc-bright text-sm font-mono">
              深度 {{ peak.depth.toFixed(2) }} mm
            </span>
            <span class="text-osc-amber text-xs font-mono">
              {{ peak.amplitude.toFixed(3) }}
            </span>
          </div>
          <div class="flex items-center gap-3 text-osc-muted text-[11px] font-mono mt-0.5">
            <span>t={{ peak.time.toFixed(2) }} µs</span>
            <span>显著度 {{ peak.prominence.toFixed(3) }}</span>
          </div>
        </div>
        <MapPin
          :class="[
            'w-4 h-4 flex-shrink-0 transition-opacity',
            highlightIndex === idx ? 'text-osc-red opacity-100' : 'text-osc-muted opacity-0 group-hover:opacity-50',
          ]"
        />
      </div>
    </div>
  </div>
</template>
