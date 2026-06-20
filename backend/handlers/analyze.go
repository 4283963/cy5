package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ultrasonic-backend/services"
)

type AnalyzeHandler struct {
	pythonRunner *services.PythonRunner
	historyStore *services.HistoryStore
}

func NewAnalyzeHandler(pr *services.PythonRunner, hs *services.HistoryStore) *AnalyzeHandler {
	return &AnalyzeHandler{
		pythonRunner: pr,
		historyStore: hs,
	}
}

type AnalyzeParams struct {
	GaussianSigma  float64 `json:"gaussianSigma"`
	FFTWindow      string  `json:"fftWindow"`
	PeakProminence float64 `json:"peakProminence"`
	PeakDistance   int     `json:"peakDistance"`
}

type AnalyzeRequest struct {
	Signal     []float64     `json:"signal"`
	SampleRate float64       `json:"sampleRate"`
	FileName   string        `json:"fileName"`
	Params     AnalyzeParams `json:"params"`
}

func (h *AnalyzeHandler) Analyze(c *gin.Context) {
	var req AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if len(req.Signal) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signal is empty"})
		return
	}

	sampleRate := req.SampleRate
	if sampleRate <= 0 {
		sampleRate = 100.0
	}

	params := req.Params
	if params.GaussianSigma <= 0 {
		params.GaussianSigma = 3
	}
	if params.FFTWindow == "" {
		params.FFTWindow = "hann"
	}
	if params.PeakProminence <= 0 {
		params.PeakProminence = 0.2
	}
	if params.PeakDistance <= 0 {
		params.PeakDistance = 10
	}

	pyInput := services.PythonInput{
		Signal:     req.Signal,
		SampleRate: sampleRate,
		Params: services.PythonInputParams{
			GaussianSigma:  params.GaussianSigma,
			FFTWindow:      params.FFTWindow,
			PeakProminence: params.PeakProminence,
			PeakDistance:   params.PeakDistance,
		},
	}

	pyOutput, err := h.pythonRunner.Run(pyInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "signal processing failed: " + err.Error()})
		return
	}

	id := services.GenerateID()
	fileName := req.FileName
	if fileName == "" {
		fileName = "signal.json"
	}

	peaks := make([]services.Peak, len(pyOutput.Peaks))
	for i, p := range pyOutput.Peaks {
		peaks[i] = services.Peak{
			Index:      p.Index,
			Time:       p.Time,
			Depth:      p.Depth,
			Amplitude:  p.Amplitude,
			Prominence: p.Prominence,
		}
	}

	record := &services.AnalysisRecord{
		ID:        id,
		FileName:  fileName,
		CreatedAt: time.Now().Format(time.RFC3339),
		Raw:       req.Signal,
		Filtered:  pyOutput.Filtered,
		FFTFreq:   pyOutput.FFTFreq,
		FFTMag:    pyOutput.FFTMag,
		Peaks:     peaks,
		Stats: services.Stats{
			Points:         pyOutput.Stats.Points,
			SNRImprovement: pyOutput.Stats.SNRImprovement,
			PeakCount:      pyOutput.Stats.PeakCount,
			MaxAmplitude:   pyOutput.Stats.MaxAmplitude,
		},
	}

	if err := h.historyStore.Save(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record: " + err.Error()})
		return
	}

	if err := h.historyStore.SaveUploadCSV(id, req.Signal); err != nil {
	}

	c.JSON(http.StatusOK, record)
}
