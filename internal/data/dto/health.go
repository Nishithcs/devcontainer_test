package dto

type HealthResponse struct {
	Status      string     `json:"status"`
	Timestamp   int64      `json:"timestamp"`
	Environment string     `json:"environment"`
	System      SystemInfo `json:"system"`
}

type SystemInfo struct {
	Version      string `json:"version"`
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"goroutines"`
	MemoryUsage  string `json:"memory_usage"`
	NumCPU       int    `json:"num_cpu"`
	Uptime       string `json:"uptime"`
}
