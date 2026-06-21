<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, MarkLineComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { Layers, MapPin, Activity, Gauge } from 'lucide-vue-next'
import type { SegmentPeak } from '@/types'

use([LineChart, GridComponent, TooltipComponent, MarkLineComponent, CanvasRenderer])

const props = defineProps<{
  segments: SegmentPeak[]
}>()

const segmentColors = [
  { main: '#00e676', glow: 'rgba(0, 230, 118, 0.4)' },
  { main: '#00e5ff', glow: 'rgba(0, 229, 255, 0.4)' },
  { main: '#ffb300', glow: 'rgba(255, 179, 0, 0.4)' },
  { main: '#e040fb', glow: 'rgba(224, 64, 251, 0.4)' },
  { main: '#ff3d3d', glow: 'rgba(255, 61, 61, 0.4)' },
]

const getColor = (idx: number) => segmentColors[idx % segmentColors.length]

const chartOptions = computed(() => {
  return props.segments.map((seg) => {
    const color = getColor(seg.segmentIndex)
    return {
      backgroundColor: 'transparent',
      grid: { left: 0, right: 0, top: 8, bottom: 0 },
      tooltip: {
        show: false,
      },
      xAxis: {
        type: 'category',
        show: false,
        data: seg.timeAxis.map((t) => t.toFixed(2)),
      },
      yAxis: {
        type: 'value',
        show: false,
      },
      series: [
        {
          type: 'line',
          data: seg.waveform,
          showSymbol: false,
          smooth: true,
          lineStyle: {
            color: color.main,
            width: 1.5,
            shadowColor: color.glow,
            shadowBlur: 3,
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0, y: 0, x2: 0, y2: 1,
              colorStops: [
                { offset: 0, color: color.main + '30' },
                { offset: 1, color: color.main + '00' },
              ],
            },
          },
          markLine: {
            silent: true,
            symbol: 'none',
            lineStyle: { color: color.main, type: 'dashed', width: 1 },
            label: { show: false },
            data: [{ xAxis: seg.peakTime.toFixed(2) }],
          },
        },
      ],
    }
  })
})
</script>

<template>
  <div class="space-y-3">
    <div class="flex items-center gap-2 mb-1">
      <Layers class="w-4 h-4 text-osc-green" />
      <h4 class="text-osc-bright text-sm font-medium">
        波形分段分析
      </h4>
      <span class="text-osc-muted text-xs">
        平均切 {{ segments.length }} 段 · 每段最尖锐峰标记
      </span>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3">
      <div
        v-for="(seg, idx) in segments"
        :key="seg.segmentIndex"
        class="bg-osc-card border rounded-lg overflow-hidden transition-all duration-200 hover:scale-[1.02] hover:shadow-card"
        :style="{
          borderColor: getColor(idx).main + '40',
          boxShadow: '0 0 12px ' + getColor(idx).glow + '20',
        }"
      >
        <div
          class="px-3 py-2 flex items-center justify-between"
          :style="{ background: getColor(idx).main + '10' }"
        >
          <div class="flex items-center gap-2">
            <div
              class="w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold"
              :style="{ background: getColor(idx).main + '30', color: getColor(idx).main }"
            >
              {{ idx + 1 }}
            </div>
            <span class="text-osc-bright text-xs font-mono">
              {{ seg.startTime.toFixed(1) }}–{{ seg.endTime.toFixed(1) }} µs
            </span>
          </div>
          <MapPin class="w-3 h-3" :style="{ color: getColor(idx).main }" />
        </div>

        <div class="px-2 py-1 h-[70px]">
          <v-chart :option="chartOptions[idx]" autoresize class="w-full h-full" />
        </div>

        <div class="px-3 py-2 border-t border-osc-border/50 space-y-1.5">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-1 text-osc-muted text-[10px]">
              <Activity class="w-3 h-3" />
              <span>深度</span>
            </div>
            <span class="text-xs font-mono font-medium" :style="{ color: getColor(idx).main }">
              {{ seg.peakDepth.toFixed(2) }} mm
            </span>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center gap-1 text-osc-muted text-[10px]">
              <MapPin class="w-3 h-3" />
              <span>位置</span>
            </div>
            <span class="text-xs font-mono text-osc-bright">
              {{ seg.peakTime.toFixed(2) }} µs
            </span>
          </div>

          <div class="flex items-center justify-between">
            <div class="flex items-center gap-1 text-osc-muted text-[10px]">
              <Gauge class="w-3 h-3" />
              <span>锐度</span>
            </div>
            <span class="text-xs font-mono text-osc-amber">
              {{ seg.sharpness.toFixed(2) }}
            </span>
          </div>

          <div class="pt-1 border-t border-osc-border/50">
            <div class="text-osc-muted text-[10px] mb-0.5">幅值</div>
            <div class="text-sm font-mono font-bold" :style="{ color: getColor(idx).main }">
              {{ seg.amplitude.toFixed(4) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
