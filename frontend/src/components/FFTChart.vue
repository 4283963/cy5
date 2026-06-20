<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, MarkLineComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([BarChart, LineChart, GridComponent, TooltipComponent, MarkLineComponent, CanvasRenderer])

const props = defineProps<{
  freq: number[]
  mag: number[]
}>()

const chartOption = computed(() => {
  const maxIdx = props.mag.indexOf(Math.max(...props.mag))
  const dominantFreq = maxIdx >= 0 ? props.freq[maxIdx] : 0

  return {
    backgroundColor: 'transparent',
    grid: {
      left: 50,
      right: 20,
      top: 20,
      bottom: 35,
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 20, 25, 0.95)',
      borderColor: '#21262d',
      textStyle: { color: '#c9d1d9', fontSize: 11 },
      formatter: (params: any) => {
        const p = params[0]
        return `<div style="font-family:monospace;font-size:11px;">
          <strong>${p.axisValue.toFixed(3)} MHz</strong><br/>
          幅值: ${p.value.toFixed(5)}
        </div>`
      },
    },
    xAxis: {
      type: 'category',
      data: props.freq.map((f) => +f.toFixed(2)),
      axisLine: { lineStyle: { color: '#21262d' } },
      axisLabel: {
        color: '#6e7681',
        fontSize: 10,
        fontFamily: 'monospace',
        interval: 'auto',
      },
      splitLine: { show: false },
      name: '频率 (MHz)',
      nameTextStyle: { color: '#6e7681', fontSize: 11, padding: [6, 0, 0, 0] },
    },
    yAxis: {
      type: 'value',
      name: '幅值',
      nameTextStyle: { color: '#6e7681', fontSize: 11 },
      axisLine: { show: false },
      axisLabel: { color: '#6e7681', fontSize: 10, fontFamily: 'monospace' },
      splitLine: {
        lineStyle: { color: 'rgba(0, 229, 255, 0.06)', type: 'dashed' },
      },
    },
    series: [
      {
        type: 'bar',
        data: props.mag,
        barWidth: '60%',
        itemStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(0, 229, 255, 0.9)' },
              { offset: 1, color: 'rgba(0, 229, 255, 0.15)' },
            ],
          },
        },
        emphasis: {
          itemStyle: { color: '#00e5ff' },
        },
        markLine: {
          silent: true,
          symbol: 'none',
          lineStyle: {
            color: '#ffb300',
            type: 'dashed',
            width: 1,
          },
          label: {
            formatter: `主频 ${dominantFreq.toFixed(2)} MHz`,
            color: '#ffb300',
            fontSize: 10,
            fontFamily: 'monospace',
          },
          data: [{ xAxis: dominantFreq.toFixed(2) }],
        },
      },
    ],
  }
})
</script>

<template>
  <div class="w-full h-full min-h-[200px]">
    <v-chart class="w-full h-full" :option="chartOption" autoresize />
  </div>
</template>
