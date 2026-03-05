terraform {
  required_providers {
    matillion-streaming = {
      source = "matillion-public/matillion-streaming"
    }
  }
}

provider "matillion-streaming" {
  # Account ID for your Matillion Data Productivity Cloud account
  account_id = "your-account-id"

  # Region where your account is hosted (either "eu" or "us")
  region = "eu"

  # Authentication credentials are provided via environment variables:
  # MATILLION_CLIENT_ID
  # MATILLION_CLIENT_SECRET
}
