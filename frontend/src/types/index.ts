export interface Peak {
  index: number
  time: number
  depth: number
  amplitude: number
  prominence: number
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
  raw: number[]
  filtered: number[]
  fftFreq: number[]
  fftMag: number[]
  peaks: Peak[]
  stats: Stats
}

export interface AnalyzeParams {
  gaussianSigma: number
  fftWindow: string
  peakProminence: number
  peakDistance: number
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
