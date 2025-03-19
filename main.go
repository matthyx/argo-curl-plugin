package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestBody struct {
	// Add any request body fields if needed
}

type ResponseBody struct {
	Output struct {
		Parameters []map[string]interface{} `json:"parameters"`
	} `json:"output"`
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

	// Decode the request body if necessary
	/*
	   var reqBody RequestBody
	   err := json.NewDecoder(r.Body).Decode(&reqBody)
	   if err != nil {
	           http.Error(w, err.Error(), http.StatusBadRequest)
	           return
	   }
	*/

	// Generate the response data
	params := []map[string]interface{}{
		{
			"name":  "param1",
			"value": "value1",
		},
		{
			"name":  "param2",
			"value": 123,
		},
		{
			"name":  "param3",
			"value": true,
		},
	}

	response := ResponseBody{}
	response.Output.Parameters = params

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
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
