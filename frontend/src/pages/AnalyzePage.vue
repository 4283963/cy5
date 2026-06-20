<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Zap } from 'lucide-vue-next'
import OscilloCard from '@/components/OscilloCard.vue'
import SignalUploader from '@/components/SignalUploader.vue'
import ParamsPanel from '@/components/ParamsPanel.vue'
import WaveformChart from '@/components/WaveformChart.vue'
import FFTChart from '@/components/FFTChart.vue'
import PeakList from '@/components/PeakList.vue'
import StatsCards from '@/components/StatsCards.vue'
import { analyzeSignal, getSample } from '@/api'
import type { AnalysisResult, AnalyzeParams, Peak, Stats } from '@/types'

const signal = ref<number[]>([])
const fileName = ref('')
const sampleRate = ref(100)
const loading = ref(false)
const result = ref<AnalysisResult | null>(null)
const highlightPeakIndex = ref<number>(-1)

const params = ref<AnalyzeParams>({
  gaussianSigma: 3,
  fftWindow: 'hann',
  peakProminence: 0.2,
  peakDistance: 10,
})

const hasSignal = computed(() => signal.value.length > 0)

const displayRaw = computed(() => result.value?.raw ?? signal.value)
const displayFiltered = computed(() => result.value?.filtered ?? [])
const displayPeaks = computed<Peak[]>(() => result.value?.peaks ?? [])
const displayStats = computed<Stats | null>(() => result.value?.stats ?? null)
const displayFFTFreq = computed(() => result.value?.fftFreq ?? [])
const displayFFTMag = computed(() => result.value?.fftMag ?? [])

const handleSignalLoaded = (payload: { signal: number[]; fileName: string }) => {
  signal.value = payload.signal
  fileName.value = payload.fileName
  result.value = null
  highlightPeakIndex.value = -1
}

const loadSample = async () => {
  loading.value = true
  try {
    const sample = await getSample()
    signal.value = sample.signal
    sampleRate.value = sample.sampleRate
    fileName.value = sample.fileName
    result.value = null
    highlightPeakIndex.value = -1
  } finally {
    loading.value = false
  }
}

const runAnalysis = async () => {
  if (!hasSignal.value) return
  loading.value = true
  try {
    const res = await analyzeSignal(
      signal.value,
      sampleRate.value,
      fileName.value || 'signal.json',
      params.value,
    )
    result.value = res
  } catch (e: any) {
    console.error('Analysis failed:', e)
    alert('分析失败: ' + (e.response?.data?.error || e.message))
  } finally {
    loading.value = false
  }
}

const handlePeakHighlight = (idx: number) => {
  highlightPeakIndex.value = highlightPeakIndex.value === idx ? -1 : idx
}

onMounted(() => {
  loadSample()
})
</script>

<template>
  <div class="min-h-[calc(100vh-57px)] bg-osc-bg grid-bg-fine">
    <div class="max-w-[1600px] mx-auto px-6 py-5">
      <div class="flex items-center justify-between mb-5">
        <div>
          <h2 class="text-osc-bright text-xl font-display font-bold flex items-center gap-2">
            <Zap class="w-5 h-5 text-osc-green" />
            信号分析工作台
          </h2>
          <p class="text-osc-muted text-sm mt-1">
            上传超声波回波原始信号，通过高斯滤波与傅里叶变换提取缺陷特征峰
          </p>
        </div>
        <div v-if="result" class="text-right">
          <div class="text-osc-muted text-xs">分析记录 ID</div>
          <div class="text-osc-green font-mono text-sm">{{ result.id }}</div>
        </div>
      </div>

      <StatsCards :stats="displayStats" class="mb-5" />

      <div class="grid grid-cols-12 gap-5">
        <div class="col-span-12 lg:col-span-3 space-y-5">
          <OscilloCard title="信号输入">
            <SignalUploader @signal-loaded="handleSignalLoaded" />
          </OscilloCard>

          <OscilloCard title="参数配置">
            <ParamsPanel
              v-model="params"
              :sample-rate="sampleRate"
              :loading="loading"
              :has-signal="hasSignal"
              @run="runAnalysis"
              @load-sample="loadSample"
            />
          </OscilloCard>

          <OscilloCard title="缺陷峰列表" :subtitle="`共 ${displayPeaks.length} 个`">
            <PeakList
              :peaks="displayPeaks"
              :highlight-index="highlightPeakIndex"
              @highlight="handlePeakHighlight"
            />
          </OscilloCard>
        </div>

        <div class="col-span-12 lg:col-span-9 space-y-5">
          <OscilloCard title="波形对比图" subtitle="原始信号 vs 高斯滤波后信号" class="h-[420px]">
            <div class="h-[360px]">
              <WaveformChart
                v-if="hasSignal"
                :raw-signal="displayRaw"
                :filtered-signal="displayFiltered"
                :sample-rate="sampleRate"
                :peaks="displayPeaks"
                :highlight-index="highlightPeakIndex"
              />
              <div v-else class="h-full flex items-center justify-center text-osc-muted text-sm">
                请先上传或加载信号
              </div>
            </div>
          </OscilloCard>

          <OscilloCard title="FFT 频谱分析" subtitle="频域特征分布" class="h-[300px]">
            <div class="h-[240px]">
              <FFTChart v-if="result" :freq="displayFFTFreq" :mag="displayFFTMag" />
              <div v-else class="h-full flex items-center justify-center text-osc-muted text-sm">
                运行分析后显示频谱
              </div>
            </div>
          </OscilloCard>
        </div>
      </div>
    </div>
  </div>
</template>
