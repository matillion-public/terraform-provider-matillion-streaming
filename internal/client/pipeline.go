package client

import "fmt"

// PipelineService provides CRUD operations for streaming pipelines
type PipelineService struct {
	client *Client
}

// PipelineOptions defines the parameters for creating or updating a streaming pipeline
type PipelineOptions struct {
	Name               string            // Required: Name of the streaming pipeline
	AgentId            string            // Required: ID of the agent to use
	StreamingSource    interface{}       // Required: Source configuration
	StreamingTarget    interface{}       // Required: Target configuration
	AdvancedProperties map[string]string // Optional: Advanced properties for configuration
}

// pipelineCreate defines the JSON structure for API requests
type pipelineCreate struct {
	Name               string            `json:"name"`
	AgentId            string            `json:"agentId"`
	StreamingSource    interface{}       `json:"streamingSource"`
	StreamingTarget    interface{}       `json:"streamingTarget"`
	AdvancedProperties map[string]string `json:"advancedProperties,omitempty"`
}

// PipelineListResponse represents the paginated response from the list streaming pipelines endpoint
type PipelineListResponse struct {
	Results []Pipeline `json:"results"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Total   int        `json:"total"`
}

// List retrieves all streaming pipelines for a project across all pages
func (s *PipelineService) List(projectId string) ([]Pipeline, error) {
	var allPipelines []Pipeline
	page := 0

	for {
		pipelinesUrl := fmt.Sprintf("%s/projects/%s/streaming-pipelines?page=%d", s.client.apiURL, projectId, page)

		var response PipelineListResponse
		_, err := s.client.doRequest("GET", pipelinesUrl, nil, StatusOK, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Results) == 0 {
			break
		}

		allPipelines = append(allPipelines, response.Results...)

		// Stop if we got less than a full page (last page reached)
		// This may fetch one extra empty page if total items is exactly divisible by page size,
		// but the simplicity is worth the negligible cost of one extra API call
		if len(response.Results) < response.Size {
			break
		}

		page++
	}

	return allPipelines, nil
}

// Create submits a new streaming pipeline to the API and returns its ID
func (s *PipelineService) Create(projectId string, opts PipelineOptions) (string, error) {
	pipelineUrl := fmt.Sprintf("%s/projects/%s/streaming-pipelines", s.client.apiURL, projectId)

	pipeline := pipelineCreate{
		Name:               opts.Name,
		AgentId:            opts.AgentId,
		StreamingSource:    opts.StreamingSource,
		StreamingTarget:    opts.StreamingTarget,
		AdvancedProperties: opts.AdvancedProperties,
	}

	var created Pipeline
	_, err := s.client.doRequest("POST", pipelineUrl, pipeline, StatusCreated, &created)
	if err != nil {
		return "", err
	}

	return created.StreamingPipelineId, nil
}

// Get retrieves a streaming pipeline by its ID from the API
func (s *PipelineService) Get(projectId string, pipelineId string) (Pipeline, error) {
	pipelineUrl := fmt.Sprintf("%s/projects/%s/streaming-pipelines/%s", s.client.apiURL, projectId, pipelineId)

	var pipeline Pipeline
	_, err := s.client.doRequest("GET", pipelineUrl, nil, StatusOK, &pipeline)
	if err != nil {
		return Pipeline{}, err
	}

	return pipeline, nil
}

// Replace updates an existing streaming pipeline with new configuration
func (s *PipelineService) Replace(projectId string, pipelineId string, opts PipelineOptions) (string, error) {
	pipelineUrl := fmt.Sprintf("%s/projects/%s/streaming-pipelines/%s", s.client.apiURL, projectId, pipelineId)

	pipeline := pipelineCreate{
		Name:               opts.Name,
		AgentId:            opts.AgentId,
		StreamingSource:    opts.StreamingSource,
		StreamingTarget:    opts.StreamingTarget,
		AdvancedProperties: opts.AdvancedProperties,
	}

	var created Pipeline
	_, err := s.client.doRequest("PUT", pipelineUrl, pipeline, StatusOK, &created)
	if err != nil {
		return "", err
	}

	return created.StreamingPipelineId, nil
}

// Delete removes a streaming pipeline from the system
func (s *PipelineService) Delete(projectId string, pipelineId string) error {
	pipelineUrl := fmt.Sprintf("%s/projects/%s/streaming-pipelines/%s", s.client.apiURL, projectId, pipelineId)
	_, err := s.client.doRequest("DELETE", pipelineUrl, nil, StatusNoContent, nil)
	return err
}
