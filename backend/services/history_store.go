package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Peak struct {
	Index      int     `json:"index"`
	Time       float64 `json:"time"`
	Depth      float64 `json:"depth"`
	Amplitude  float64 `json:"amplitude"`
	Prominence float64 `json:"prominence"`
}

type Stats struct {
	Points         int     `json:"points"`
	SNRImprovement float64 `json:"snrImprovement"`
	PeakCount      int     `json:"peakCount"`
	MaxAmplitude   float64 `json:"maxAmplitude"`
}

type AnalysisRecord struct {
	ID         string        `json:"id"`
	FileName   string        `json:"fileName"`
	CreatedAt  string        `json:"createdAt"`
	SampleRate int           `json:"sampleRate"`
	Raw        []float64     `json:"raw"`
	Filtered   []float64     `json:"filtered"`
	FFTFreq    []float64     `json:"fftFreq"`
	FFTMag     []float64     `json:"fftMag"`
	Peaks      []Peak        `json:"peaks"`
	Segments   []SegmentPeak `json:"segments"`
	Stats      Stats         `json:"stats"`
}

type HistorySummary struct {
	ID           string  `json:"id"`
	FileName     string  `json:"fileName"`
	CreatedAt    string  `json:"createdAt"`
	PeakCount    int     `json:"peakCount"`
	Points       int     `json:"points"`
	MaxAmplitude float64 `json:"maxAmplitude"`
}

type HistoryStore struct {
	historyDir string
	uploadsDir string
}

func NewHistoryStore(historyDir, uploadsDir string) *HistoryStore {
	return &HistoryStore{
		historyDir: historyDir,
		uploadsDir: uploadsDir,
	}
}

func GenerateID() string {
	now := time.Now()
	return now.Format("20060102-150405") + fmt.Sprintf("-%03d", now.Nanosecond()/1e6)
}

func (hs *HistoryStore) Save(record *AnalysisRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(hs.historyDir, record.ID+".json")
	return os.WriteFile(path, data, 0644)
}

func (hs *HistoryStore) SaveUploadCSV(id string, signal []float64) error {
	path := filepath.Join(hs.uploadsDir, id+".csv")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range signal {
		if _, err := fmt.Fprintf(f, "%.6f\n", v); err != nil {
			return err
		}
	}
	return nil
}

func (hs *HistoryStore) Get(id string) (*AnalysisRecord, error) {
	path := filepath.Join(hs.historyDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var record AnalysisRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (hs *HistoryStore) List() ([]HistorySummary, error) {
	entries, err := os.ReadDir(hs.historyDir)
	if err != nil {
		return nil, err
	}
	var summaries []HistorySummary
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		path := filepath.Join(hs.historyDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var record AnalysisRecord
		if err := json.Unmarshal(data, &record); err != nil {
			continue
		}
		summaries = append(summaries, HistorySummary{
			ID:           record.ID,
			FileName:     record.FileName,
			CreatedAt:    record.CreatedAt,
			PeakCount:    record.Stats.PeakCount,
			Points:       record.Stats.Points,
			MaxAmplitude: record.Stats.MaxAmplitude,
		})
	}
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].CreatedAt > summaries[j].CreatedAt
	})
	return summaries, nil
}

func (hs *HistoryStore) Delete(id string) error {
	historyPath := filepath.Join(hs.historyDir, id+".json")
	if err := os.Remove(historyPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	uploadPath := filepath.Join(hs.uploadsDir, id+".csv")
	if err := os.Remove(uploadPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
