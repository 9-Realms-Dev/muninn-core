package munnincore

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

// Reader functions
func ReadHttpFile(path string) ([]HttpRequest, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Store all the requests in a file
	var requests []HttpRequest

	// Keep track of current request
	var req HttpRequest

	// Body variables
	var isBody bool
	var bodyLines []string

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line starts with "###"
		if strings.HasPrefix(line, "###") {
			// Check if we have a current request to add
			if req.Title != "" {
				req.Body = bytes.NewBufferString(strings.Join(bodyLines, "\n"))
				requests = append(requests, req)
			}

			// Create a new request
			req = HttpRequest{
				Title:   strings.TrimSpace(line[3:]),
				Headers: make(map[string]string),
			}
			isBody = false
			bodyLines = nil
		}

		// Check if the line starts with an HTTP method
		if isMethodLine(line) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				req.Method = parts[0]
				req.URL = parts[1]
			}
		}

		// Read through all the headers
		if strings.Contains(line, ":") && !isMethodLine(line) && !isBody {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Headers[key] = value
			}
		}

		// Read through the body
		if line == "" {
			isBody = true
		}

		if isBody && line != "" {
			bodyLines = append(bodyLines, line)
		}

		// TODO: Add post call script details
	}

	// Append the last test data to the tests slice
	if req.Title != "" {
		req.Body = bytes.NewBufferString(strings.Join(bodyLines, "\n"))
		requests = append(requests, req)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func isMethodLine(line string) bool {
	return strings.HasPrefix(line, "GET") ||
		strings.HasPrefix(line, "POST") ||
		strings.HasPrefix(line, "PUT") ||
		strings.HasPrefix(line, "DELETE")
}

// Send Functions
func SendHttpRequests(requests []HttpRequest) ([]HttpResponse, error) {
	if len(requests) == 0 {
		return nil, errors.New("no requets found in provided file")
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout per request
	}

	var resps []HttpResponse

	for _, request := range requests {
		req, err := http.NewRequest(request.Method, request.URL, request.Body)
		if err != nil {
			// TODO: Add some form of error handling for collecting multiple errors to return
			continue
		}

		// set headers
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)

		resps = append(resps, HttpResponse{
			Response: resp,
			Error:    err,
		})
	}

	return resps, nil
}

func SendHttpRequest(request HttpRequest) (*HttpResponse, error) {
	client := &http.Client{
		Timeout: 10 * time.Second, // Timeout per request
	}

	req, err := http.NewRequest(request.Method, request.URL, request.Body)
	if err != nil {
		return nil, err
	}

	// set headers
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)

	return &HttpResponse{
		Response: resp,
		Error:    err,
	}, nil
}
