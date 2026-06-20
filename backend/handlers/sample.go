package handlers

import (
	"math"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SampleHandler struct{}

func NewSampleHandler() *SampleHandler {
	return &SampleHandler{}
}

type SampleResponse struct {
	Signal     []float64 `json:"signal"`
	SampleRate float64   `json:"sampleRate"`
	FileName   string    `json:"fileName"`
}

func (h *SampleHandler) GetSample(c *gin.Context) {
	n := 2048
	sampleRate := 100.0
	signal := make([]float64, n)

	rng := rand.New(rand.NewSource(42))

	baseline := 0.02
	for i := 0; i < n; i++ {
		signal[i] = (rng.Float64() - 0.5) * 2.0 * baseline
	}

	addGaussianEcho(signal, sampleRate, 5.0, 0.9, 0.3)
	addGaussianEcho(signal, sampleRate, 12.0, 0.6, 0.5)
	addGaussianEcho(signal, sampleRate, 2.0, 1.2, 0.15)

	for i := 0; i < n; i++ {
		signal[i] += (rng.Float64() - 0.5) * 2.0 * 0.03
	}

	c.JSON(http.StatusOK, SampleResponse{
		Signal:     signal,
		SampleRate: sampleRate,
		FileName:   "sample-steel.json",
	})
}

func addGaussianEcho(signal []float64, sampleRate, centerUs, amplitude, widthUs float64) {
	n := len(signal)
	centerIdx := centerUs * sampleRate
	sigmaIdx := widthUs * sampleRate

	for i := 0; i < n; i++ {
		dx := float64(i) - centerIdx
		signal[i] += amplitude * math.Exp(-dx*dx/(2.0*sigmaIdx*sigmaIdx))
	}
}
