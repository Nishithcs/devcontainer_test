package jobs

import (
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// DevpodOutputParser contains regex patterns for parsing devpod command output
type DevpodOutputParser struct {
	// Machine creation patterns
	machineCreatePattern *regexp.Regexp

	// URL and port patterns
	urlPattern *regexp.Regexp

	// Success patterns
	successStopPattern   *regexp.Regexp
	successDeletePattern *regexp.Regexp

	// ANSI escape codes for log cleaning
	ansiEscapePattern *regexp.Regexp

	// Log parsing pattern
	logPattern *regexp.Regexp
}

// NewDevpodOutputParser creates a new parser with compiled regex patterns
func NewDevpodOutputParser() *DevpodOutputParser {
	return &DevpodOutputParser{
		machineCreatePattern: regexp.MustCompile(`Create machine '([^']+)'`),
		urlPattern:           regexp.MustCompile(`Successfully opened (http[^\s]+)`),
		successStopPattern:   regexp.MustCompile(`Successfully stopped`),
		successDeletePattern: regexp.MustCompile(`Successfully deleted workspace`),
		ansiEscapePattern:    regexp.MustCompile(`\x1b\[[0-9;]*m`),
		logPattern:           regexp.MustCompile(`^(?P<time>\d{2}:\d{2}:\d{2})\s+(?P<type>\w+)\s+(?P<text>.+)$`),
	}
}

// ExtractMachineName extracts machine name from devpod output
func (p *DevpodOutputParser) ExtractMachineName(message string) (string, bool) {
	matches := p.machineCreatePattern.FindStringSubmatch(message)
	if len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

// ExtractURL extracts workspace URL from devpod output
func (p *DevpodOutputParser) ExtractURL(message string) (string, bool) {
	matches := p.urlPattern.FindStringSubmatch(message)
	if len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

// ExtractPortFromURL extracts port number from a URL
func (p *DevpodOutputParser) ExtractPortFromURL(rawURL string) int {
	// Parse the URL to extract the port
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Failed to parse URL: %v\n", err)
	}

	// Extract port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		log.Printf("Unexpected host format: %s\n", parsedURL.Host)
	}

	portStr := hostParts[1]
	internalPort, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("Invalid port: %s\n", portStr)
	}

	return internalPort

}

// IsSuccessStop checks if the message indicates successful stop
func (p *DevpodOutputParser) IsSuccessStop(message string) bool {
	return p.successStopPattern.MatchString(message)
}

// IsSuccessDelete checks if the message indicates successful deletion
func (p *DevpodOutputParser) IsSuccessDelete(message string) bool {
	return p.successDeletePattern.MatchString(message)
}

// ParseLogMessage parses a log message and returns time, type, and text
func (p *DevpodOutputParser) ParseLogMessage(log string) (time, logType, text string) {
	log = strings.TrimSpace(log)

	// Remove ANSI color codes
	log = p.ansiEscapePattern.ReplaceAllString(log, "")

	// Match against clean log
	match := p.logPattern.FindStringSubmatch(log)
	if match != nil {
		return match[1], match[2], match[3]
	}
	return "", "", ""
}

// Global parser instance
var devpodParser = NewDevpodOutputParser()

// Convenience functions for backward compatibility
func getClearLog(log string) (string, string, string) {
	return devpodParser.ParseLogMessage(log)
}
