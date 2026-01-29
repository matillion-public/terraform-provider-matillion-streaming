resource "matillion-streaming_agent" "example" {
  name           = "my-streaming-agent"
  description    = "Example streaming agent"
  deployment     = "fargate"
  cloud_provider = "aws"
}

# Output the computed agent_id
output "agent_id" {
  value       = matillion-streaming_agent.example.agent_id
  description = "The unique identifier of the created agent"
}
