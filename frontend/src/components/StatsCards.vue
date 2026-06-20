<script setup lang="ts">
import { Activity, TrendingUp, AlertTriangle, Gauge } from 'lucide-vue-next'
import type { Stats } from '@/types'

defineProps<{
  stats: Stats | null
}>()
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
    <div class="bg-osc-card border border-osc-border rounded-lg p-3 relative overflow-hidden">
      <div class="absolute top-2 right-2">
        <Activity class="w-4 h-4 text-osc-green/40" />
      </div>
      <div class="text-osc-muted text-xs mb-1">采样点数</div>
      <div class="text-osc-bright text-xl font-display font-bold text-glow-green">
        {{ stats?.points?.toLocaleString() || '--' }}
      </div>
      <div class="text-osc-muted text-[10px] mt-0.5">原始信号采样</div>
    </div>

    <div class="bg-osc-card border border-osc-border rounded-lg p-3 relative overflow-hidden">
      <div class="absolute top-2 right-2">
        <TrendingUp class="w-4 h-4 text-osc-cyan/40" />
      </div>
      <div class="text-osc-muted text-xs mb-1">信噪比提升</div>
      <div class="text-osc-cyan text-xl font-display font-bold text-glow-cyan">
        {{ stats ? stats.snrImprovement.toFixed(1) : '--' }}
        <span class="text-sm font-normal text-osc-muted">dB</span>
      </div>
      <div class="text-osc-muted text-[10px] mt-0.5">滤波后噪声抑制</div>
    </div>

    <div class="bg-osc-card border border-osc-border rounded-lg p-3 relative overflow-hidden">
      <div class="absolute top-2 right-2">
        <AlertTriangle class="w-4 h-4 text-osc-red/40" />
      </div>
      <div class="text-osc-muted text-xs mb-1">检出缺陷峰</div>
      <div class="text-osc-red text-xl font-display font-bold text-glow-red">
        {{ stats?.peakCount ?? '--' }}
        <span class="text-sm font-normal text-osc-muted">个</span>
      </div>
      <div class="text-osc-muted text-[10px] mt-0.5">潜在缺陷位置</div>
    </div>

    <div class="bg-osc-card border border-osc-border rounded-lg p-3 relative overflow-hidden">
      <div class="absolute top-2 right-2">
        <Gauge class="w-4 h-4 text-osc-amber/40" />
      </div>
      <div class="text-osc-muted text-xs mb-1">最大幅值</div>
      <div class="text-osc-amber text-xl font-display font-bold">
        {{ stats ? stats.maxAmplitude.toFixed(3) : '--' }}
      </div>
      <div class="text-osc-muted text-[10px] mt-0.5">滤波后峰值</div>
    </div>
  </div>
</template>
