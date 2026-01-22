package client

// Agent represents an agent in the Matillion API
type Agent struct {
	AgentId       string `json:"agentId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Deployment    string `json:"deployment"`
	CloudProvider string `json:"cloudProvider"`
	AgentType     string `json:"agentType"`
}
