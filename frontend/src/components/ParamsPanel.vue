<script setup lang="ts">
import { computed } from 'vue'
import { Sliders, Play, Loader2 } from 'lucide-vue-next'
import type { AnalyzeParams } from '@/types'

const props = defineProps<{
  modelValue: AnalyzeParams
  sampleRate: number
  loading: boolean
  hasSignal: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: AnalyzeParams): void
  (e: 'update:sampleRate', value: number): void
  (e: 'run'): void
  (e: 'loadSample'): void
}>()

const params = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const updateParam = <K extends keyof AnalyzeParams>(key: K, value: AnalyzeParams[K]) => {
  emit('update:modelValue', { ...props.modelValue, [key]: value })
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <Sliders class="w-4 h-4 text-osc-green" />
        <h4 class="text-osc-bright text-sm font-medium">处理参数</h4>
      </div>
    </div>

    <div class="space-y-3">
      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="text-osc-muted text-xs">采样率 (MHz)</label>
          <span class="text-osc-green text-xs font-mono">{{ sampleRate.toFixed(1) }}</span>
        </div>
        <input
          type="number"
          :value="sampleRate"
          @input="emit('update:sampleRate', Number(($event.target as HTMLInputElement).value))"
          min="1"
          step="0.5"
          class="w-full bg-osc-bg border border-osc-border rounded-md px-3 py-1.5 text-osc-bright text-sm font-mono focus:outline-none focus:border-osc-green/50 focus:ring-1 focus:ring-osc-green/30"
        />
      </div>

      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="text-osc-muted text-xs">高斯滤波 σ</label>
          <span class="text-osc-green text-xs font-mono">{{ params.gaussianSigma.toFixed(1) }}</span>
        </div>
        <input
          type="range"
          :value="params.gaussianSigma"
          @input="updateParam('gaussianSigma', Number(($event.target as HTMLInputElement).value))"
          min="0.5"
          max="20"
          step="0.5"
          class="w-full h-1.5 bg-osc-border rounded-full appearance-none cursor-pointer accent-osc-green"
        />
      </div>

      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="text-osc-muted text-xs">FFT 窗函数</label>
        </div>
        <select
          :value="params.fftWindow"
          @change="updateParam('fftWindow', ($event.target as HTMLSelectElement).value)"
          class="w-full bg-osc-bg border border-osc-border rounded-md px-3 py-1.5 text-osc-bright text-sm font-mono focus:outline-none focus:border-osc-green/50 focus:ring-1 focus:ring-osc-green/30"
        >
          <option value="hann">Hann (汉宁窗)</option>
          <option value="hamming">Hamming (汉明窗)</option>
          <option value="rect">Rectangular (矩形窗)</option>
        </select>
      </div>

      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="text-osc-muted text-xs">峰显著度阈值</label>
          <span class="text-osc-green text-xs font-mono">{{ params.peakProminence.toFixed(2) }}</span>
        </div>
        <input
          type="range"
          :value="params.peakProminence"
          @input="updateParam('peakProminence', Number(($event.target as HTMLInputElement).value))"
          min="0.01"
          max="1"
          step="0.01"
          class="w-full h-1.5 bg-osc-border rounded-full appearance-none cursor-pointer accent-osc-green"
        />
      </div>

      <div>
        <div class="flex items-center justify-between mb-1">
          <label class="text-osc-muted text-xs">峰最小间距 (点)</label>
          <span class="text-osc-green text-xs font-mono">{{ params.peakDistance }}</span>
        </div>
        <input
          type="range"
          :value="params.peakDistance"
          @input="updateParam('peakDistance', Number(($event.target as HTMLInputElement).value))"
          min="1"
          max="100"
          step="1"
          class="w-full h-1.5 bg-osc-border rounded-full appearance-none cursor-pointer accent-osc-green"
        />
      </div>
    </div>

    <div class="pt-2 space-y-2">
      <button
        @click="emit('run')"
        :disabled="loading || !hasSignal"
        :class="[
          'w-full py-2.5 rounded-md font-medium text-sm flex items-center justify-center gap-2 transition-all duration-200',
          loading || !hasSignal
            ? 'bg-osc-border text-osc-muted cursor-not-allowed'
            : 'bg-osc-green text-osc-bg hover:shadow-green-glow active:scale-[0.98]',
        ]"
      >
        <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
        <Play v-else class="w-4 h-4" />
        {{ loading ? '分析中...' : '运行分析' }}
      </button>

      <button
        @click="emit('loadSample')"
        :disabled="loading"
        class="w-full py-2 rounded-md text-sm border border-osc-border text-osc-text hover:border-osc-cyan/50 hover:text-osc-cyan transition-all duration-200"
      >
        加载示例信号
      </button>
    </div>
  </div>
</template>
