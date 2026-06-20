import axios from 'axios'
import type { AnalysisResult, AnalyzeParams, HistoryListResponse, SampleResponse } from '@/types'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080/api',
  timeout: 60000,
})

export async function getSample(): Promise<SampleResponse> {
  const { data } = await api.get<SampleResponse>('/sample')
  return data
}

export async function analyzeSignal(
  signal: number[],
  sampleRate: number,
  fileName: string,
  params: AnalyzeParams,
): Promise<AnalysisResult> {
  const { data } = await api.post<AnalysisResult>('/analyze', {
    signal,
    sampleRate,
    fileName,
    params,
  })
  return data
}

export async function getHistoryList(): Promise<HistoryListResponse> {
  const { data } = await api.get<HistoryListResponse>('/history')
  return data
}

export async function getHistoryRecord(id: string): Promise<AnalysisResult> {
  const { data } = await api.get<AnalysisResult>(`/history/${id}`)
  return data
}

export async function deleteHistoryRecord(id: string): Promise<{ ok: boolean }> {
  const { data } = await api.delete(`/history/${id}`)
  return data
}

export default api
