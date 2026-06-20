package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type PythonRunner struct {
	Interpreter string
	ScriptPath  string
	Timeout     time.Duration
}

type PythonInputParams struct {
	GaussianSigma  float64 `json:"gaussian_sigma"`
	FFTWindow      string  `json:"fft_window"`
	PeakProminence float64 `json:"peak_prominence"`
	PeakDistance   int     `json:"peak_distance"`
}

type PythonInput struct {
	Signal     []float64         `json:"signal"`
	SampleRate float64           `json:"sample_rate"`
	Params     PythonInputParams `json:"params"`
}

type PythonPeak struct {
	Index      int     `json:"index"`
	Time       float64 `json:"time"`
	Depth      float64 `json:"depth"`
	Amplitude  float64 `json:"amplitude"`
	Prominence float64 `json:"prominence"`
}

type PythonStats struct {
	Points         int     `json:"points"`
	SNRImprovement float64 `json:"snr_improvement"`
	PeakCount      int     `json:"peak_count"`
	MaxAmplitude   float64 `json:"max_amplitude"`
}

type PythonOutput struct {
	Filtered []float64    `json:"filtered"`
	FFTFreq  []float64    `json:"fft_freq"`
	FFTMag   []float64    `json:"fft_mag"`
	Peaks    []PythonPeak `json:"peaks"`
	Stats    PythonStats  `json:"stats"`
}

func NewPythonRunner(interpreter, scriptPath string) *PythonRunner {
	return &PythonRunner{
		Interpreter: interpreter,
		ScriptPath:  scriptPath,
		Timeout:     60 * time.Second,
	}
}

func (pr *PythonRunner) Run(input PythonInput) (*PythonOutput, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pr.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pr.Interpreter, pr.ScriptPath)
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrStr := stderr.String()
		if stderrStr != "" {
			return nil, fmt.Errorf("python script failed: %w\nstderr: %s", err, stderrStr)
		}
		return nil, fmt.Errorf("python script failed: %w", err)
	}

	stdoutStr := stdout.String()
	var output PythonOutput
	if err := json.Unmarshal([]byte(stdoutStr), &output); err != nil {
		return nil, fmt.Errorf("failed to parse python output: %w\nstdout: %s", err, stdoutStr)
	}

	return &output, nil
}
