package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ultrasonic-backend/services"
)

const maxUploadSize = 50 * 1024 * 1024

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

type UploadResponse struct {
	Signal     []float64 `json:"signal"`
	SampleRate float64   `json:"sampleRate"`
	FileName   string    `json:"fileName"`
	Points     int       `json:"points"`
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传文件: " + err.Error()})
		return
	}

	if file.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，最大支持 50MB"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法读取文件: " + err.Error()})
		return
	}
	defer src.Close()

	content, err := services.ReadAll(src, maxUploadSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsed, err := services.ParseSignalFile(file.Filename, content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UploadResponse{
		Signal:     parsed.Signal,
		SampleRate: parsed.SampleRate,
		FileName:   parsed.FileName,
		Points:     len(parsed.Signal),
	})
}
