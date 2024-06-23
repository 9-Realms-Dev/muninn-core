package formats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type JsonResponse struct {
	Status     string              `json:"status"`
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       interface{}         `json:"body,omitempty"`
	RawBody    string              `json:"rawBody"`
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
		// Declare a variable to hold the parsed JSON data
		var jsonBody interface{}

		// Unmarshal the JSON response body into the jsonBody variable
		err = json.Unmarshal(body, &jsonBody)
		if err != nil {
			return JsonResponse{}, fmt.Errorf("failed to parse JSON response body: %v", err)
		}

		// Check the type of the unmarshaled data
		switch data := jsonBody.(type) {
		case map[string]interface{}:
			// If it's a single object, assign it to the Body field
			formattedResp.Body = data
		case []interface{}:
			// If it's an array, convert it to a slice of maps and assign it to the Body field
			var mapSlice []map[string]interface{}
			for _, obj := range data {
				if mapObj, ok := obj.(map[string]interface{}); ok {
					mapSlice = append(mapSlice, mapObj)
				}
			}
			formattedResp.Body = mapSlice
		default:
			// If the JSON structure is neither an object nor an array, return an error
			return JsonResponse{}, fmt.Errorf("unsupported JSON structure")
		}
	}

	return formattedResp, nil
}

func (r JsonResponse) CliRender(isRenderedBody bool) string {
	var builder strings.Builder

	// Define styles
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("green"))
	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("blue"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("gray"))

	// Render status and status code
	builder.WriteString(titleStyle.Render("Status: "))
	builder.WriteString(valueStyle.Render(r.Status))
	builder.WriteString("\n")
	builder.WriteString(titleStyle.Render("Status Code: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", r.StatusCode)))
	builder.WriteString("\n\n")

	// Render headers
	builder.WriteString(titleStyle.Render("Headers:\n"))
	for key, values := range r.Headers {
		builder.WriteString(keyStyle.Render(key + ": "))
		builder.WriteString(valueStyle.Render(strings.Join(values, ", ")))
		builder.WriteString("\n")
	}
	builder.WriteString("\n")

	// Render body
	if isRenderedBody && r.Body != nil {
		builder.WriteString(titleStyle.Render("Body:\n"))

		switch body := r.Body.(type) {
		case map[string]interface{}:
			// Render a single object
			for key, value := range body {
				builder.WriteString(keyStyle.Render(key + ": "))
				builder.WriteString(valueStyle.Render(fmt.Sprintf("%v", value)))
				builder.WriteString("\n")
			}
		case []map[string]interface{}:
			// Render an array of objects
			for i, obj := range body {
				builder.WriteString(fmt.Sprintf("Object %d:\n", i+1))
				for key, value := range obj {
					builder.WriteString(keyStyle.Render(key + ": "))
					builder.WriteString(valueStyle.Render(fmt.Sprintf("%v", value)))
					builder.WriteString("\n")
				}
				builder.WriteString("\n")
			}
		default:
			builder.WriteString(valueStyle.Render(fmt.Sprintf("%v", r.Body)))
		}

		builder.WriteString("\n")
	}

	// Render raw body
	builder.WriteString(titleStyle.Render("Raw Body:\n"))
	builder.WriteString(valueStyle.Render(r.RawBody))

	return builder.String()
}
