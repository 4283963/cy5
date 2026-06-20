package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ultrasonic-backend/handlers"
	"ultrasonic-backend/services"
)

func resolvePythonInterpreter() string {
	if v := os.Getenv("ULTRASONIC_PYTHON"); v != "" {
		return v
	}
	return filepath.Join("..", "python", ".venv", "bin", "python")
}

func resolvePythonScript() string {
	if v := os.Getenv("ULTRASONIC_PYTHON_SCRIPT"); v != "" {
		return v
	}
	return filepath.Join("..", "python", "analyze.py")
}

func ensureDataDirs() error {
	for _, dir := range []string{
		filepath.Join("data", "uploads"),
		filepath.Join("data", "history"),
	} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func limitRequestBody(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

func main() {
	if err := ensureDataDirs(); err != nil {
		log.Fatalf("failed to create data directories: %v", err)
	}

	pythonRunner := services.NewPythonRunner(
		resolvePythonInterpreter(),
		resolvePythonScript(),
	)

	historyStore := services.NewHistoryStore(
		filepath.Join("data", "history"),
		filepath.Join("data", "uploads"),
	)

	analyzeHandler := handlers.NewAnalyzeHandler(pythonRunner, historyStore)
	historyHandler := handlers.NewHistoryHandler(historyStore)
	sampleHandler := handlers.NewSampleHandler()
	uploadHandler := handlers.NewUploadHandler()

	r := gin.Default()
	r.MaxMultipartMemory = 50 << 20

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	r.Use(limitRequestBody(50 << 20))

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/upload", uploadHandler.Upload)
	r.GET("/api/sample", sampleHandler.GetSample)
	r.POST("/api/analyze", analyzeHandler.Analyze)
	r.GET("/api/history", historyHandler.ListHistory)
	r.GET("/api/history/:id", historyHandler.GetHistory)
	r.DELETE("/api/history/:id", historyHandler.DeleteHistory)

	log.Println("🚀 Ultrasonic NDT backend starting on :8080")
	log.Printf("   Python interpreter: %s", pythonRunner.Interpreter)
	log.Printf("   Python script:      %s", pythonRunner.ScriptPath)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
