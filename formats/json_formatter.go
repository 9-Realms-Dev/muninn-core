package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type JsonResponse struct {
	Status     string                 `json:"status"`
	StatusCode int                    `json:"statusCode"`
	Headers    map[string][]string    `json:"headers"`
	Body       map[string]interface{} `json:"body,omitempty"`
	RawBody    string                 `json:"rawBody"`
}

func FormatJSONResponse(resp *http.Response) (JsonResponse, error) {
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return JsonResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Create a formatted response
	formattedResp := JsonResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		RawBody:    string(body),
	}

	// Check if the response has a JSON content type
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// Parse the JSON response body
		var jsonBody map[string]interface{}
		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			return JsonResponse{}, fmt.Errorf("failed to parse JSON response body: %v", err)
		}
		formattedResp.Body = jsonBody
	}

	return formattedResp, nil
}
