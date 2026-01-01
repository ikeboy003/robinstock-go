package robinstock_go

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ikeboy003/robinstock-go/models"
	"github.com/ikeboy003/robinstock-go/utils"
)

const (
	defaultTimeout = 10 * time.Second
)

var (
	ErrNotAuthenticated = errors.New("not authenticated")
	ErrInvalidResponse  = errors.New("invalid response from API")
)

// Client represents a Robinhood API client.
type Client struct {
	httpClient        *http.Client
	phoenixHTTPClient *http.Client
	auth              *models.Auth
}

// NewClient creates a new Robinhood API client.
func NewClient() *Client {
	// Standard client for most endpoints
	standardClient := &http.Client{
		Timeout: defaultTimeout,
	}

	// Custom client for Phoenix endpoint with TLS 1.2 only and specific cipher suites
	phoenixTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			MaxVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
			ServerName: "phoenix.robinhood.com",
		},
	}
	phoenixClient := &http.Client{
		Timeout:   defaultTimeout,
		Transport: phoenixTransport,
	}

	return &Client{
		httpClient:        standardClient,
		phoenixHTTPClient: phoenixClient,
	}
}

// SetTimeout sets the HTTP request timeout.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetAuth sets authentication credentials.
func (c *Client) SetAuth(auth *models.Auth) {
	c.auth = auth
}

 // GetAuth returns the current authentication credentials.
func (c *Client) GetAuth() *models.Auth {
	return c.auth
}

// IsAuthenticated returns true if the client is authenticated.
func (c *Client) IsAuthenticated() bool {
	return c.auth != nil && c.auth.AccessToken != ""
}

// doRequest executes an HTTP request with proper headers.
func (c *Client) doRequest(ctx context.Context, method, urlStr string, body interface{}, authenticated bool) (*models.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set standard headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip,deflate,br")
	req.Header.Set("Accept-Language", "en-US,en;q=1")
	req.Header.Set("X-Robinhood-API-Version", models.ApiVersion)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "robinstock_go/1.0")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add authentication if required
	if authenticated {
		if !c.IsAuthenticated() {
			return nil, ErrNotAuthenticated
		}
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.auth.TokenType, c.auth.AccessToken))
	}

	// Use Phoenix client for phoenix.robinhood.com endpoints
	httpClient := c.httpClient
	if strings.Contains(urlStr, "phoenix.robinhood.com") {
		httpClient = c.phoenixHTTPClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	return parseResponse(resp)
}

// Get executes a GET request.
func (c *Client) Get(ctx context.Context, urlStr string, params map[string]string, authenticated bool) (*models.Response, error) {
	if params != nil {
		urlStr = utils.BuildURL(urlStr, params)
	}
	return c.doRequest(ctx, http.MethodGet, urlStr, nil, authenticated)
}

// Post executes a POST request.
func (c *Client) Post(ctx context.Context, urlStr string, body interface{}, authenticated bool) (*models.Response, error) {
	return c.doRequest(ctx, http.MethodPost, urlStr, body, authenticated)
}

func parseResponse(resp *http.Response) (*models.Response, error) {
	// Read body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// Decompress if needed
	var reader io.Reader = bytes.NewReader(bodyBytes)
	encoding := resp.Header.Get("Content-Encoding")
	switch encoding {
	case "gzip":
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("gzip decompress: %w", err)
		}
	case "deflate":
		reader, err = zlib.NewReader(reader)
		if err != nil {
			return nil, fmt.Errorf("deflate decompress: %w", err)
		}
	}

	// Decode JSON
	var data map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	response := &models.Response{
		StatusCode: resp.StatusCode,
		Data:       data,
	}

	// Extract results if present (paginated response)
	if results, ok := data["results"].([]interface{}); ok {
		for _, r := range results {
			if m, ok := r.(map[string]interface{}); ok {
				response.Results = append(response.Results, m)
			}
		}
	}

	return response, nil
}

// FetchAllPages fetches all pages of a paginated response.
func (c *Client) FetchAllPages(ctx context.Context, initialURL string, authenticated bool) ([]map[string]interface{}, error) {
	var allResults []map[string]interface{}
	nextURL := initialURL

	for nextURL != "" {
		resp, err := c.Get(ctx, nextURL, nil, authenticated)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, resp.Results...)

		// Check for next page
		if next, ok := resp.Data["next"].(string); ok && next != "" {
			nextURL = next
		} else {
			nextURL = ""
		}
	}

	return allResults, nil
}

