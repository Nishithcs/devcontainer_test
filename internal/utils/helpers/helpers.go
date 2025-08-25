package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func GetEnvironment() string {
	env := gin.Mode()
	if env == gin.ReleaseMode {
		return "production"
	}
	if env == gin.TestMode {
		return "test"
	}
	return "development"
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	if days > 0 {
		return fmt.Sprintf("%dd%dh%dm%ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func LoadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	if err := godotenv.Load(); err != nil {
		if env == "development" {
			fmt.Printf("Error loading .env file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("No .env file found, using environment variables")
	} else {
		fmt.Println("Loaded .env file")
	}
}

func PascalToSnakeCase(s string) string {
	var result []string
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, "_")
		}
		result = append(result, string(unicode.ToLower(r)))
	}
	return strings.Join(result, "")
}

func extractURL(line string) string {
	// crude but effective: look for "http..." in the line
	parts := strings.Fields(line)
	for _, part := range parts {
		if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
			return part
		}
	}
	return ""
}
