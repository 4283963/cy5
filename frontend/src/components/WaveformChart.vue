<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, MarkPointComponent, LegendComponent, DataZoomComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import type { Peak } from '@/types'

use([LineChart, GridComponent, TooltipComponent, MarkPointComponent, LegendComponent, DataZoomComponent, CanvasRenderer])

const props = defineProps<{
  rawSignal: number[]
  filteredSignal: number[]
  sampleRate: number
  peaks: Peak[]
  highlightIndex?: number
}>()

const emit = defineEmits<{
  (e: 'peakClick', index: number): void
}>()

const timeAxis = computed(() => {
  const n = props.rawSignal.length
  const dt = 1 / props.sampleRate
  return Array.from({ length: n }, (_, i) => +(i * dt).toFixed(3))
})

const peakMarkData = computed(() => {
  return props.peaks.map((p, idx) => ({
    name: `缺陷峰${idx + 1}`,
    coord: [p.time, p.amplitude],
    value: `#${idx + 1}`,
    itemStyle: {
      color: props.highlightIndex === idx ? '#ff3d3d' : 'rgba(255, 61, 61, 0.85)',
    },
    label: {
      show: true,
      color: '#fff',
      fontSize: 10,
      fontWeight: 'bold',
    },
  }))
})

const chartOption = computed(() => ({
  backgroundColor: 'transparent',
  grid: {
    left: 50,
    right: 20,
    top: 30,
    bottom: 40,
  },
  legend: {
    data: ['原始信号', '滤波后信号'],
    textStyle: { color: '#c9d1d9', fontSize: 11 },
    top: 0,
    right: 10,
    itemWidth: 14,
    itemHeight: 2,
  },
  tooltip: {
    trigger: 'axis',
    backgroundColor: 'rgba(15, 20, 25, 0.95)',
    borderColor: '#21262d',
    textStyle: { color: '#c9d1d9', fontSize: 11 },
    formatter: (params: any) => {
      const t = params[0]?.axisValue
      let html = `<div class="font-mono text-xs"><strong>时间: ${t} µs</strong></div>`
      params.forEach((p: any) => {
        const color = p.color
        html += `<div style="color:${color};font-family:monospace;font-size:11px;margin-top:2px;">${p.seriesName}: ${p.value.toFixed(4)}</div>`
      })
      return html
    },
  },
  xAxis: {
    type: 'category',
    data: timeAxis.value,
    axisLine: { lineStyle: { color: '#21262d' } },
    axisLabel: {
      color: '#6e7681',
      fontSize: 10,
      fontFamily: 'monospace',
      formatter: (v: string) => `${v}`,
      interval: 'auto',
    },
    splitLine: {
      show: true,
      lineStyle: { color: 'rgba(0, 230, 118, 0.05)', type: 'dashed' },
    },
    name: '时间 (µs)',
    nameTextStyle: { color: '#6e7681', fontSize: 11, padding: [8, 0, 0, 0] },
  },
  yAxis: {
    type: 'value',
    name: '幅值',
    nameTextStyle: { color: '#6e7681', fontSize: 11 },
    axisLine: { show: false },
    axisLabel: { color: '#6e7681', fontSize: 10, fontFamily: 'monospace' },
    splitLine: {
      lineStyle: { color: 'rgba(0, 230, 118, 0.06)', type: 'dashed' },
    },
  },
  dataZoom: [
    {
      type: 'inside',
      xAxisIndex: 0,
      start: 0,
      end: 100,
    },
    {
      type: 'slider',
      xAxisIndex: 0,
      start: 0,
      end: 100,
      height: 20,
      bottom: 5,
      borderColor: '#21262d',
      backgroundColor: '#0f1419',
      fillerColor: 'rgba(0, 230, 118, 0.15)',
      handleStyle: { color: '#00e676' },
      textStyle: { color: '#6e7681', fontSize: 9 },
    },
  ],
  series: [
    {
      name: '原始信号',
      type: 'line',
      data: props.rawSignal,
      showSymbol: false,
      lineStyle: {
        color: '#00e5ff',
        width: 1,
        opacity: 0.6,
      },
      smooth: false,
      sampling: 'lttb',
    },
    {
      name: '滤波后信号',
      type: 'line',
      data: props.filteredSignal,
      showSymbol: false,
      lineStyle: {
        color: '#00e676',
        width: 1.5,
        shadowColor: 'rgba(0, 230, 118, 0.5)',
        shadowBlur: 4,
      },
      smooth: false,
      sampling: 'lttb',
      markPoint: {
        symbol: 'pin',
        symbolSize: 30,
        data: peakMarkData.value,
        emphasis: { scale: 1.3 },
      },
    },
  ],
}))
</script>

<template>
  <div class="w-full h-full min-h-[320px]">
    <v-chart class="w-full h-full" :option="chartOption" autoresize />
  </div>
</template>
