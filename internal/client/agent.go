package client

import "fmt"

// AgentService handles operations on agents
type AgentService struct {
	client *Client
}

// CreateAgentOptions contains all options for creating an agent
type CreateAgentOptions struct {
	Name          string // Name of the agent
	Description   string // Description of the agent
	Deployment    string // Deployment type (e.g., "fargate")
	CloudProvider string // Cloud provider (e.g., "aws")
	AgentType     string // Type of agent (e.g., "streaming")
}

// UpdateAgentOptions contains all options for updating an agent
type UpdateAgentOptions struct {
	Name        string // New name for the agent
	Description string // New description for the agent
}

// AgentListResponse represents the paginated response from the list agents endpoint
type AgentListResponse struct {
	Results []Agent `json:"results"`
	Page    int     `json:"page"`
	Size    int     `json:"size"`
	Total   int     `json:"total"`
}

// List retrieves all agents across all pages
func (s *AgentService) List() ([]Agent, error) {
	var allAgents []Agent
	page := 0

	for {
		agentUrl := fmt.Sprintf("%s/agents?page=%d", s.client.apiURL, page)

		var response AgentListResponse
		_, err := s.client.doRequest("GET", agentUrl, nil, StatusOK, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Results) == 0 {
			break
		}

		allAgents = append(allAgents, response.Results...)

		// Stop if we got less than a full page (last page reached)
		// This may fetch one extra empty page if total items is exactly divisible by page size,
		// but the simplicity is worth the negligible cost of one extra API call
		if len(response.Results) < response.Size {
			break
		}

		page++
	}

	return allAgents, nil
}

// Create creates a new agent using the options pattern
func (s *AgentService) Create(opts CreateAgentOptions) (string, error) {
	agentUrl := fmt.Sprintf("%s/agents", s.client.apiURL)

	requestBody := struct {
		Name          string `json:"name"`
		Description   string `json:"description,omitempty"`
		Deployment    string `json:"deployment"`
		CloudProvider string `json:"cloudProvider"`
		AgentType     string `json:"agentType"`
	}{
		Name:          opts.Name,
		Description:   opts.Description,
		Deployment:    opts.Deployment,
		CloudProvider: opts.CloudProvider,
		AgentType:     opts.AgentType,
	}

	var created struct {
		AgentId string `json:"agentId"`
	}
	_, err := s.client.doRequest("POST", agentUrl, requestBody, StatusCreated, &created)
	if err != nil {
		return "", err
	}

	return created.AgentId, nil
}

// Get retrieves an agent by its ID
func (s *AgentService) Get(agentId string) (Agent, error) {
	agentUrl := fmt.Sprintf("%s/agents/%s", s.client.apiURL, agentId)

	var agent Agent
	_, err := s.client.doRequest("GET", agentUrl, nil, StatusOK, &agent)
	if err != nil {
		return Agent{}, err
	}

	return agent, nil
}

// Update updates an existing agent using the options pattern
func (s *AgentService) Update(agentId string, opts UpdateAgentOptions) error {
	agentUrl := fmt.Sprintf("%s/agents/%s", s.client.apiURL, agentId)

	requestBody := struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		Name:        opts.Name,
		Description: opts.Description,
	}

	_, err := s.client.doRequest("PATCH", agentUrl, requestBody, StatusAccepted, nil)
	return err
}

// Delete deletes an agent
func (s *AgentService) Delete(agentId string) error {
	agentUrl := fmt.Sprintf("%s/agents/%s", s.client.apiURL, agentId)
	_, err := s.client.doRequest("DELETE", agentUrl, nil, StatusAccepted, nil)
	return err
}

// GetCredentials retrieves the credentials for an agent
func (s *AgentService) GetCredentials(agentId string) (AgentCredentials, error) {
	agentUrl := fmt.Sprintf("%s/agents/%s/credentials", s.client.apiURL, agentId)

	var credentials AgentCredentials
	_, err := s.client.doRequest("GET", agentUrl, nil, StatusOK, &credentials)
	if err != nil {
		return AgentCredentials{}, err
	}

	return credentials, nil
}
