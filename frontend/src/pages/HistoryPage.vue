<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { History, Trash2, Eye, Clock, Activity, AlertTriangle, ArrowLeft } from 'lucide-vue-next'
import OscilloCard from '@/components/OscilloCard.vue'
import WaveformChart from '@/components/WaveformChart.vue'
import FFTChart from '@/components/FFTChart.vue'
import PeakList from '@/components/PeakList.vue'
import StatsCards from '@/components/StatsCards.vue'
import { getHistoryList, getHistoryRecord, deleteHistoryRecord } from '@/api'
import type { HistorySummary, AnalysisResult } from '@/types'

const records = ref<HistorySummary[]>([])
const loading = ref(false)
const detail = ref<AnalysisResult | null>(null)
const detailLoading = ref(false)
const highlightPeakIndex = ref<number>(-1)

const loadList = async () => {
  loading.value = true
  try {
    const res = await getHistoryList()
    records.value = res.records
  } finally {
    loading.value = false
  }
}

const viewDetail = async (id: string) => {
  detailLoading.value = true
  try {
    detail.value = await getHistoryRecord(id)
    highlightPeakIndex.value = -1
  } finally {
    detailLoading.value = false
  }
}

const closeDetail = () => {
  detail.value = null
}

const deleteRecord = async (id: string, e: Event) => {
  e.stopPropagation()
  if (!confirm('确定删除这条记录吗？')) return
  try {
    await deleteHistoryRecord(id)
    if (detail.value?.id === id) {
      detail.value = null
    }
    loadList()
  } catch (e: any) {
    alert('删除失败: ' + (e.response?.data?.error || e.message))
  }
}

const formatDate = (iso: string) => {
  const d = new Date(iso)
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

const handlePeakHighlight = (idx: number) => {
  highlightPeakIndex.value = highlightPeakIndex.value === idx ? -1 : idx
}

onMounted(() => {
  loadList()
})
</script>

<template>
  <div class="min-h-[calc(100vh-57px)] bg-osc-bg grid-bg-fine">
    <div class="max-w-[1400px] mx-auto px-6 py-5">
      <div class="flex items-center justify-between mb-5">
        <div>
          <h2 class="text-osc-bright text-xl font-display font-bold flex items-center gap-2">
            <History class="w-5 h-5 text-osc-green" />
            检测历史
          </h2>
          <p class="text-osc-muted text-sm mt-1">
            过往分析记录回顾与管理
          </p>
        </div>
        <button
          v-if="detail"
          @click="closeDetail"
          class="flex items-center gap-2 px-4 py-2 text-sm text-osc-text border border-osc-border rounded-md hover:border-osc-cyan/50 hover:text-osc-cyan transition-all"
        >
          <ArrowLeft class="w-4 h-4" />
          返回列表
        </button>
      </div>

      <div v-if="!detail">
        <OscilloCard>
          <template #title>
            <span class="flex items-center gap-2">
              <Clock class="w-4 h-4 text-osc-green" />
              分析记录
            </span>
          </template>
          <template v-if="loading">
            <div class="py-12 text-center text-osc-muted">加载中...</div>
          </template>
          <template v-else-if="records.length === 0">
            <div class="py-12 text-center text-osc-muted">
              暂无历史记录
              <div class="text-xs mt-2">去工作台运行一次分析吧</div>
            </div>
          </template>
          <template v-else>
            <div class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-osc-border text-osc-muted text-xs uppercase tracking-wider">
                    <th class="text-left py-2.5 px-3 font-medium">文件名</th>
                    <th class="text-left py-2.5 px-3 font-medium">分析时间</th>
                    <th class="text-right py-2.5 px-3 font-medium">采样点</th>
                    <th class="text-right py-2.5 px-3 font-medium">缺陷峰</th>
                    <th class="text-right py-2.5 px-3 font-medium">最大幅值</th>
                    <th class="text-right py-2.5 px-3 font-medium">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="record in records"
                    :key="record.id"
                    class="border-b border-osc-border/50 hover:bg-osc-green/[0.03] transition-colors cursor-pointer"
                    @click="viewDetail(record.id)"
                  >
                    <td class="py-2.5 px-3">
                      <div class="text-osc-bright font-mono text-sm truncate max-w-[240px]">
                        {{ record.fileName }}
                      </div>
                      <div class="text-osc-muted text-[11px] font-mono">{{ record.id }}</div>
                    </td>
                    <td class="py-2.5 px-3 text-osc-muted text-sm">
                      {{ formatDate(record.createdAt) }}
                    </td>
                    <td class="py-2.5 px-3 text-right text-osc-text font-mono">
                      {{ record.points.toLocaleString() }}
                    </td>
                    <td class="py-2.5 px-3 text-right">
                      <span
                        :class="[
                          'inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-mono',
                          record.peakCount > 0
                            ? 'bg-osc-red/15 text-osc-red'
                            : 'bg-osc-green/15 text-osc-green',
                        ]"
                      >
                        <AlertTriangle v-if="record.peakCount > 0" class="w-3 h-3" />
                        <Activity v-else class="w-3 h-3" />
                        {{ record.peakCount }}
                      </span>
                    </td>
                    <td class="py-2.5 px-3 text-right text-osc-amber font-mono">
                      {{ record.maxAmplitude.toFixed(3) }}
                    </td>
                    <td class="py-2.5 px-3 text-right">
                      <div class="flex items-center justify-end gap-1">
                        <button
                          @click="viewDetail(record.id)"
                          class="p-1.5 text-osc-muted hover:text-osc-cyan transition-colors"
                          title="查看详情"
                        >
                          <Eye class="w-4 h-4" />
                        </button>
                        <button
                          @click="(e) => deleteRecord(record.id, e)"
                          class="p-1.5 text-osc-muted hover:text-osc-red transition-colors"
                          title="删除"
                        >
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </template>
        </OscilloCard>
      </div>

      <div v-else class="space-y-5">
        <OscilloCard>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-osc-bright font-mono text-base">{{ detail.fileName }}</div>
              <div class="text-osc-muted text-xs font-mono mt-0.5">
                ID: {{ detail.id }} · {{ formatDate(detail.createdAt) }}
              </div>
            </div>
          </div>
        </OscilloCard>

        <StatsCards :stats="detail.stats" />

        <div class="grid grid-cols-12 gap-5">
          <div class="col-span-12 lg:col-span-3">
            <OscilloCard title="缺陷峰列表" :subtitle="`共 ${detail.peaks.length} 个`">
              <PeakList
                :peaks="detail.peaks"
                :highlight-index="highlightPeakIndex"
                @highlight="handlePeakHighlight"
              />
            </OscilloCard>
          </div>
          <div class="col-span-12 lg:col-span-9 space-y-5">
            <OscilloCard title="波形对比图" class="h-[420px]">
              <div class="h-[360px]">
                <WaveformChart
                  :raw-signal="detail.raw"
                  :filtered-signal="detail.filtered"
                  :sample-rate="100"
                  :peaks="detail.peaks"
                  :highlight-index="highlightPeakIndex"
                />
              </div>
            </OscilloCard>
            <OscilloCard title="FFT 频谱分析" class="h-[260px]">
              <div class="h-[200px]">
                <FFTChart :freq="detail.fftFreq" :mag="detail.fftMag" />
              </div>
            </OscilloCard>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
