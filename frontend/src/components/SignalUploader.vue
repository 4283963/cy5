<script setup lang="ts">
import { ref } from 'vue'
import { Upload, FileText, AlertCircle, Loader2 } from 'lucide-vue-next'
import { uploadSignalFile } from '@/api'

const emit = defineEmits<{
  (e: 'signalLoaded', payload: { signal: number[]; fileName: string; sampleRate: number }): void
}>()

const isDragging = ref(false)
const uploading = ref(false)
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
  target.value = ''
}

const processFile = async (file: File) => {
  error.value = ''
  uploading.value = true
  try {
    const result = await uploadSignalFile(file)
    fileInfo.value = {
      name: result.fileName,
      size: file.size,
      points: result.points,
    }
    emit('signalLoaded', {
      signal: result.signal,
      fileName: result.fileName,
      sampleRate: result.sampleRate,
    })
  } catch (e: any) {
    fileInfo.value = null
    error.value = e.response?.data?.error || e.message || '文件上传失败'
  } finally {
    uploading.value = false
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
        uploading ? 'opacity-60 pointer-events-none' : '',
        isDragging
          ? 'border-osc-green bg-osc-green/5 border-glow-green'
          : 'border-osc-border hover:border-osc-green/50 hover:bg-osc-green/[0.02]',
      ]"
    >
      <input
        type="file"
        class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
        accept=".csv,.json,.txt,.dat"
        :disabled="uploading"
        @change="handleFileSelect"
      />
      <div class="flex flex-col items-center gap-2 pointer-events-none">
        <div class="w-10 h-10 rounded-full bg-osc-green/10 flex items-center justify-center">
          <Loader2 v-if="uploading" class="w-5 h-5 text-osc-green animate-spin" />
          <Upload v-else class="w-5 h-5 text-osc-green" />
        </div>
        <div class="text-osc-bright text-sm font-medium">
          {{ uploading ? '上传解析中...' : '拖拽文件到此处，或点击选择' }}
        </div>
        <div class="text-osc-muted text-xs">
          支持 CSV / JSON / TXT / DAT 格式，最大 50MB
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
      <div class="w-2 h-2 rounded-full bg-osc-green shadow-green-glow"></div>
    </div>

    <div v-if="error" class="bg-osc-red/10 border border-osc-red/30 rounded-md px-3 py-2 flex items-center gap-2">
      <AlertCircle class="w-4 h-4 text-osc-red flex-shrink-0" />
      <span class="text-osc-red text-sm">{{ error }}</span>
    </div>
  </div>
</template>
