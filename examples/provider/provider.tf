terraform {
  required_providers {
    clickup = {
      source = "catdevman/clickup"
    }
  }
}

provider "clickup" {
  api_token = "API_KEY"
}
