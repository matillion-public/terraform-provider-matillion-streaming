#!/bin/bash
# Example: Import an existing streaming pipeline into Terraform state

# The import ID format is: project_id:pipeline_id
# Where:
#   - project_id is your Matillion Data Productivity Cloud project ID
#   - pipeline_id is the existing pipeline's unique identifier

# Import syntax:
terraform import matillion-streaming_pipeline.example "your-project-id:existing-pipeline-id"

# Example with actual values:
# terraform import matillion-streaming_pipeline.production_pipeline "proj-abc123:pipe-xyz789"

# After importing, you can view the imported configuration:
# terraform show

# Then create a matching resource block in your .tf file to manage it:
# resource "matillion-streaming_pipeline" "production_pipeline" {
#   name       = "production-pipeline"
#   project_id = "proj-abc123"
#   agent_id   = "agent-123"
#   # ... rest of configuration
# }
