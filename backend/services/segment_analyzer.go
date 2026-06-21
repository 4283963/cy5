package services

import "math"

type SegmentPeak struct {
	SegmentIndex int       `json:"segmentIndex"`
	StartIndex   int       `json:"startIndex"`
	EndIndex     int       `json:"endIndex"`
	StartTime    float64   `json:"startTime"`
	EndTime      float64   `json:"endTime"`
	PeakIndex    int       `json:"peakIndex"`
	PeakTime     float64   `json:"peakTime"`
	PeakDepth    float64   `json:"peakDepth"`
	Amplitude    float64   `json:"amplitude"`
	Sharpness    float64   `json:"sharpness"`
	Waveform     []float64 `json:"waveform"`
	TimeAxis     []float64 `json:"timeAxis"`
}

const STEEL_VELOCITY_MM_PER_US = 5.9

func ComputeSegments(filtered []float64, sampleRate float64, numSegments int, waveformPoints int) []SegmentPeak {
	if len(filtered) == 0 || numSegments <= 0 {
		return []SegmentPeak{}
	}

	if sampleRate <= 0 {
		sampleRate = 100.0
	}
	if waveformPoints <= 0 {
		waveformPoints = 200
	}

	n := len(filtered)
	segSize := n / numSegments
	if segSize < 1 {
		segSize = 1
	}

	segments := make([]SegmentPeak, 0, numSegments)

	for s := 0; s < numSegments; s++ {
		startIdx := s * segSize
		endIdx := startIdx + segSize
		if s == numSegments-1 {
			endIdx = n
		}
		if startIdx >= n {
			break
		}
		if endIdx > n {
			endIdx = n
		}

		segLen := endIdx - startIdx
		if segLen == 0 {
			continue
		}

		peakAbsIdx := startIdx
		peakAbsVal := math.Abs(filtered[startIdx])
		for i := startIdx + 1; i < endIdx; i++ {
			absVal := math.Abs(filtered[i])
			if absVal > peakAbsVal {
				peakAbsVal = absVal
				peakAbsIdx = i
			}
		}

		peakVal := filtered[peakAbsIdx]
		peakTime := float64(peakAbsIdx) / sampleRate
		sharpness := computeSharpness(filtered, peakAbsIdx, sampleRate)

		half := waveformPoints / 2
		wfStart := peakAbsIdx - half
		wfEnd := peakAbsIdx + half
		if wfStart < 0 {
			wfStart = 0
		}
		if wfEnd > n {
			wfEnd = n
		}

		wfLen := wfEnd - wfStart
		waveform := make([]float64, wfLen)
		timeAxis := make([]float64, wfLen)
		for i := 0; i < wfLen; i++ {
			idx := wfStart + i
			waveform[i] = filtered[idx]
			timeAxis[i] = float64(idx) / sampleRate
		}

		segments = append(segments, SegmentPeak{
			SegmentIndex: s,
			StartIndex:   startIdx,
			EndIndex:     endIdx,
			StartTime:    float64(startIdx) / sampleRate,
			EndTime:      float64(endIdx) / sampleRate,
			PeakIndex:    peakAbsIdx,
			PeakTime:     peakTime,
			PeakDepth:    peakTime * STEEL_VELOCITY_MM_PER_US / 2.0,
			Amplitude:    peakVal,
			Sharpness:    sharpness,
			Waveform:     waveform,
			TimeAxis:     timeAxis,
		})
	}

	return segments
}

func computeSharpness(signal []float64, peakIdx int, sampleRate float64) float64 {
	n := len(signal)
	window := int(math.Max(3, math.Min(20, float64(n)/100)))

	leftIdx := peakIdx - window
	if leftIdx < 0 {
		leftIdx = 0
	}
	rightIdx := peakIdx + window
	if rightIdx >= n {
		rightIdx = n - 1
	}

	peakVal := math.Abs(signal[peakIdx])
	var sum float64
	count := 0
	for i := leftIdx; i <= rightIdx; i++ {
		if i == peakIdx {
			continue
		}
		sum += math.Abs(signal[i])
		count++
	}
	if count == 0 {
		return 0
	}
	avg := sum / float64(count)
	if avg <= 0 {
		return 0
	}
	return peakVal / avg
}
