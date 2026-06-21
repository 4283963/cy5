export interface Peak {
  index: number
  time: number
  depth: number
  amplitude: number
  prominence: number
}

export interface SegmentPeak {
  segmentIndex: number
  startIndex: number
  endIndex: number
  startTime: number
  endTime: number
  peakIndex: number
  peakTime: number
  peakDepth: number
  amplitude: number
  sharpness: number
  waveform: number[]
  timeAxis: number[]
}

export interface Stats {
  points: number
  snrImprovement: number
  peakCount: number
  maxAmplitude: number
}

export interface AnalysisResult {
  id: string
  fileName: string
  createdAt: string
  sampleRate: number
  raw: number[]
  filtered: number[]
  fftFreq: number[]
  fftMag: number[]
  peaks: Peak[]
  segments: SegmentPeak[]
  stats: Stats
}

export interface AnalyzeParams {
  gaussianSigma: number
  fftWindow: string
  peakProminence: number
  peakDistance: number
}

export interface UploadResponse {
  signal: number[]
  sampleRate: number
  fileName: string
  points: number
}

export interface HistorySummary {
  id: string
  fileName: string
  createdAt: string
  peakCount: number
  points: number
  maxAmplitude: number
}

export interface HistoryListResponse {
  records: HistorySummary[]
}

export interface SampleResponse {
  signal: number[]
  sampleRate: number
  fileName: string
}
