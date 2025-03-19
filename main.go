package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestBody struct {
	Input struct {
		Parameters struct {
			DiscoveryURL string `json:"discoveryURL"`
		} `json:"parameters"`
	} `json:"input"`
}

type ResponseBody struct {
	Output struct {
		Parameters json.RawMessage `json:"parameters"`
	} `json:"output"`
}

type ServiceDiscoveryResponse struct {
	Version  string          `json:"version"`
	Response json.RawMessage `json:"response"`
}

func getParamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/api/v1/getparams.execute" {
		http.NotFound(w, r)
		return
	}

	// Decode the request body
	var reqBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch data from discovery API
	discoveryURL := reqBody.Input.Parameters.DiscoveryURL
	if discoveryURL == "" {
		http.Error(w, "discoveryURL is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(discoveryURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data from discovery API: %v", err), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("discovery API returned non-OK status: %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading discovery API response: %v", err), http.StatusInternalServerError)
		return
	}

	var discoveryResponse ServiceDiscoveryResponse
	err = json.Unmarshal(body, &discoveryResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling discovery API response: %v", err), http.StatusInternalServerError)
		return
	}

	response := ResponseBody{}
	response.Output.Parameters = discoveryResponse.Response

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api/v1/getparams.execute", getParamsHandler)

	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
