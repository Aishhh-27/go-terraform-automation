terraform {
  required_providers {
    local = {
      source = "hashicorp/local"
    }
  }
}

provider "local" {}

resource "local_file" "example" {
  filename = "${path.module}/hello.txt"
  content  = "Hello from DEV environment"
}
