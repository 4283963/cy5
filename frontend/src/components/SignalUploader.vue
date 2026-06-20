<script setup lang="ts">
import { ref } from 'vue'
import { Upload, FileText, AlertCircle } from 'lucide-vue-next'

const emit = defineEmits<{
  (e: 'signalLoaded', payload: { signal: number[]; fileName: string }): void
}>()

const isDragging = ref(false)
const error = ref('')
const fileInfo = ref<{ name: string; size: number; points: number } | null>(null)

const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = true
}

const handleDragLeave = () => {
  isDragging.value = false
}

const handleDrop = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    processFile(files[0])
  }
}

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = target.files
  if (files && files.length > 0) {
    processFile(files[0])
  }
}

const processFile = async (file: File) => {
  error.value = ''
  try {
    const text = await file.text()
    let signal: number[] = []

    if (file.name.toLowerCase().endsWith('.json')) {
      const data = JSON.parse(text)
      if (Array.isArray(data)) {
        signal = data.map(Number)
      } else if (data.signal && Array.isArray(data.signal)) {
        signal = data.signal.map(Number)
      } else {
        throw new Error('JSON 文件格式不正确，需要包含 signal 数组或本身是数组')
      }
    } else if (file.name.toLowerCase().endsWith('.csv') || file.name.toLowerCase().endsWith('.txt')) {
      const lines = text.trim().split(/\r?\n/)
      signal = lines
        .filter((l) => l.trim() !== '')
        .map((l) => {
          const parts = l.split(/[,;\t]/)
          return Number(parts[parts.length - 1].trim())
        })
        .filter((v) => !isNaN(v))
    } else {
      throw new Error('仅支持 .csv、.json、.txt 格式')
    }

    if (signal.length === 0) {
      throw new Error('未找到有效的信号数据')
    }

    fileInfo.value = {
      name: file.name,
      size: file.size,
      points: signal.length,
    }

    emit('signalLoaded', { signal, fileName: file.name })
  } catch (e: any) {
    error.value = e.message || '文件解析失败'
  }
}
</script>

<template>
  <div class="space-y-3">
    <div
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      @drop="handleDrop"
      :class="[
        'relative border-2 border-dashed rounded-lg p-6 text-center transition-all duration-200 cursor-pointer',
        isDragging
          ? 'border-osc-green bg-osc-green/5 border-glow-green'
          : 'border-osc-border hover:border-osc-green/50 hover:bg-osc-green/[0.02]',
      ]"
    >
      <input
        type="file"
        class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        accept=".csv,.json,.txt"
        @change="handleFileSelect"
      />
      <div class="flex flex-col items-center gap-2 pointer-events-none">
        <div class="w-10 h-10 rounded-full bg-osc-green/10 flex items-center justify-center">
          <Upload class="w-5 h-5 text-osc-green" />
        </div>
        <div class="text-osc-bright text-sm font-medium">
          拖拽文件到此处，或点击选择
        </div>
        <div class="text-osc-muted text-xs">
          支持 CSV / JSON / TXT 格式
        </div>
      </div>
    </div>

    <div v-if="fileInfo" class="bg-osc-green/5 border border-osc-green/30 rounded-md px-3 py-2 flex items-center gap-3">
      <FileText class="w-4 h-4 text-osc-green flex-shrink-0" />
      <div class="flex-1 min-w-0">
        <div class="text-osc-bright text-sm font-mono truncate">{{ fileInfo.name }}</div>
        <div class="text-osc-muted text-xs">
          {{ fileInfo.points.toLocaleString() }} 个采样点 · {{ (fileInfo.size / 1024).toFixed(1) }} KB
        </div>
      </div>
    </div>

    <div v-if="error" class="bg-osc-red/10 border border-osc-red/30 rounded-md px-3 py-2 flex items-center gap-2">
      <AlertCircle class="w-4 h-4 text-osc-red flex-shrink-0" />
      <span class="text-osc-red text-sm">{{ error }}</span>
    </div>
  </div>
</template>
