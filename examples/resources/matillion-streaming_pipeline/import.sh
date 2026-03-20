#!/bin/bash
# Example: Import an existing streaming pipeline into Terraform state

# Before importing, ensure the resource block is defined in your .tf file, for example:
# resource "matillion-streaming_pipeline" "example" {
#   name       = "production-pipeline"
#   project_id = "proj-abc123"
#   agent_id   = "agent-123"
#   # ... rest of configuration
# }

# The import ID format is: project_id:pipeline_id
terraform import matillion-streaming_pipeline.example "your-project-id:your-pipeline-id"
