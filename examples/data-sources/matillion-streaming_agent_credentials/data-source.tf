# Example: Retrieve credentials for a newly created agent
resource "matillion-streaming_agent" "example" {
  name           = "my-streaming-agent"
  description    = "Example streaming agent"
  deployment     = "fargate"
  cloud_provider = "aws"
}

data "matillion-streaming_agent_credentials" "new_agent" {
  agent_id = matillion-streaming_agent.example.agent_id
}


# Output the client_id (non-sensitive)
output "agent_client_id" {
  value       = data.matillion-streaming_agent_credentials.new_agent.client_id
  description = "The client ID for the agent"
}

# Output the client_secret (marked as sensitive)
output "agent_client_secret" {
  value       = data.matillion-streaming_agent_credentials.new_agent.client_secret
  description = "The client secret for the agent"
  sensitive   = true
}

