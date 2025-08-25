package main

import (
	"clusterix-code/internal/api_clients"
	"clusterix-code/internal/config"
	"clusterix-code/internal/data/db"
	"clusterix-code/internal/data/repositories"
	"clusterix-code/internal/services"
	"clusterix-code/internal/utils/di"
	"clusterix-code/internal/utils/helpers"
	"clusterix-code/internal/utils/logger"
	"clusterix-code/internal/utils/mongo"
	"clusterix-code/internal/utils/rabbitmq"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
)

var subdomainRegex = regexp.MustCompile(`^([a-zA-Z0-9]+)\.`)

type WorkspaceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		WorkspaceConfig struct {
			WorkerPort int `json:"worker_port"`
		} `json:"workspace_config"`
	} `json:"data"`
}

func proxyHandler(w http.ResponseWriter, r *http.Request, services *services.Services) {
	host := r.Host // e.g. abc123.example.com

	// Extract fingerprint
	matches := subdomainRegex.FindStringSubmatch(host)
	if len(matches) < 2 {
		http.Error(w, "Invalid subdomain format", http.StatusBadRequest)
		return
	}
	fingerprint := matches[1]

	// Get auth_token from query
	//auth_token := r.URL.Query().Get("auth_token")
	//if authToken == "" {
	//	http.Error(w, "Missing auth_token", http.StatusBadRequest)
	//	return
	//}

	response, err := services.Workspace.GetWorkspaceByFingerprint(r.Context(), fingerprint)

	if err != nil {
		http.Error(w, "Failed to get workspace", http.StatusInternalServerError)
		return
	}

	// Call internal API to get worker port
	//apiURL := fmt.Sprintf("http://api:8080/api/v1/workspaces/fingerprint/%s", fingerprint)
	//req, err := http.NewRequest("GET", apiURL, nil)
	//if err != nil {
	//	http.Error(w, "Failed to create API request", http.StatusInternalServerError)
	//	return
	//}
	//req.Header.Set("Authorization", "Bearer "+authToken)
	//
	//client := &http.Client{Timeout: 5 * time.Second}
	//resp, err := client.Do(req)
	//if err != nil {
	//	http.Error(w, "Failed to reach API", http.StatusBadGateway)
	//	log.Printf("API error: %v", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	bodyBytes, _ := io.ReadAll(resp.Body)
	//	log.Printf("API non-200: %s", string(bodyBytes))
	//	http.Error(w, "API error", http.StatusBadGateway)
	//	return
	//}
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	http.Error(w, "Failed to read API response", http.StatusInternalServerError)
	//	return
	//}
	//
	//var wsResp WorkspaceResponse
	//if err := json.Unmarshal(body, &wsResp); err != nil {
	//	log.Printf("JSON unmarshal error: %v", err)
	//	http.Error(w, "Invalid API response", http.StatusInternalServerError)
	//	return
	//}
	//workerPort := wsResp.Data.WorkspaceConfig.WorkerPort

	workerPort := response.WorkspaceConfig.WorkerPort
	if workerPort == nil || *workerPort == 0 {
		http.Error(w, "Invalid worker_port", http.StatusInternalServerError)
		return
	}

	worker_discovery_name := os.Getenv("WORKER_DISCOVERY_NAME")
	// Build target URL
	targetURL, err := url.Parse(fmt.Sprintf("http://%s:%d", worker_discovery_name, *workerPort))
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Rewrite request before proxying
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Set correct Host
		req.Host = targetURL.Host

		// Strip auth_token, preserve others
		query := req.URL.Query()
		query.Del("auth_token")

		// Optional: keep "folder"
		folder := query.Get("folder")
		query.Del("folder")

		req.URL.RawQuery = query.Encode()
		if folder != "" {
			if req.URL.RawQuery != "" {
				req.URL.RawQuery += "&"
			}
			req.URL.RawQuery += "folder=" + url.QueryEscape(folder)
		}
	}

	// Logging
	log.Printf("Proxying %s → %s%s?%s", host, targetURL.String(), r.URL.Path, r.URL.RawQuery)

	// Serve
	proxy.ServeHTTP(w, r)
}

func main() {
	helpers.LoadEnv()
	logger.Init(os.Getenv("APP_ENV"))
	defer logger.Sync()

	c := di.NewContainer(0)

	di.Register(c, config.Provider)
	di.Register(c, db.Provider)
	di.Register(c, mongo.Provider)
	di.Register(c, rabbitmq.Provider)
	di.Register(c, repositories.Provider)
	di.Register(c, api_clients.Provider)
	di.Register(c, services.Provider)
	c.Bootstrap()

	services := di.Make[*services.Services](c)

	port := "80"
	if p := os.Getenv("PROXY_PORT"); p != "" {
		port = p
	}

	// Inject services using closure
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxyHandler(w, r, services)
	})

	log.Printf("Reverse proxy listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
