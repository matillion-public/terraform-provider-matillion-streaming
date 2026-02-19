package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"net/url"
	"os"
)

// HTTP status code constants
const (
	StatusOK        = 200
	StatusCreated   = 201
	StatusAccepted  = 202
	StatusNoContent = 204
)

type Region string

const (
	RegionEU Region = "eu"
	RegionUS Region = "us"
)

type RegionConfig struct {
	Region      Region
	AudienceURL string
	ApiURL      string
	TokenURL    string
}

var (
	EUConfig = RegionConfig{
		Region:      RegionEU,
		AudienceURL: "https://api.matillion.com",
		ApiURL:      "https://eu1.api.matillion.com/dpc/v1",
		TokenURL:    "https://id.core.matillion.com/oauth/dpc/token",
	}

	USConfig = RegionConfig{
		Region:      RegionUS,
		AudienceURL: "https://api.matillion.com",
		ApiURL:      "https://us1.api.matillion.com/dpc/v1",
		TokenURL:    "https://id.core.matillion.com/oauth/dpc/token",
	}
)

func GetRegionConfig(region Region) RegionConfig {
	if region == RegionEU {
		return EUConfig
	} else {
		return USConfig
	}
}

type Client struct {
	accountId  string
	httpClient *http.Client
	region     Region
	apiURL     string

	Agents    *AgentService
	Pipelines *PipelineService
}

func NewClient(accountId string, region Region) (*Client, error) {
	clientId := os.Getenv("MATILLION_CLIENT_ID")
	clientSecret := os.Getenv("MATILLION_CLIENT_SECRET")

	config := GetRegionConfig(region)

	var httpClient *http.Client

	if clientId != "" && clientSecret != "" {

		oauthConfig := clientcredentials.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			TokenURL:     config.TokenURL,
			EndpointParams: url.Values{
				"audience": {config.AudienceURL},
			},
		}

		httpClient = oauthConfig.Client(context.Background())
	} else {
		return nil, fmt.Errorf("MATILLION_CLIENT_ID and MATILLION_CLIENT_SECRET environment variables are required")
	}

	client := &Client{
		accountId:  accountId,
		httpClient: httpClient,
		region:     region,
		apiURL:     config.ApiURL,
	}

	client.Agents = &AgentService{client: client}
	client.Pipelines = &PipelineService{client: client}
	return client, nil
}

// doRequest performs an HTTP request and handles common response processing
// If requestBody is provided, it will be marshaled to JSON
// If responseTarget is provided, the response will be unmarshalled into it
func (c *Client) doRequest(method, url string, requestBody interface{}, expectedStatus int, responseTarget interface{}) ([]byte, error) {
	var bodyReader io.Reader

	if requestBody != nil {
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("account-id", c.accountId)
	if requestBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		return nil, fmt.Errorf("status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	if responseTarget != nil {
		if err = json.Unmarshal(respBody, responseTarget); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return respBody, nil
}
