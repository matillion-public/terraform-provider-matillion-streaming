#!/bin/bash
# Example: Import an existing streaming agent into Terraform state

# Before importing, ensure the resource block is defined in your .tf file, for example:
# resource "matillion-streaming_agent" "example" {
#   name           = "my-streaming-agent"
#   description    = "Example streaming agent"
#   deployment     = "fargate"
#   cloud_provider = "aws"
# }

# The import ID is the agent's unique identifier
terraform import matillion-streaming_agent.example "your-agent-id"