# Terraform Backend Configuration
terraform {
  backend "gcs" {
    bucket = "dea-noctua-terraform-state"
    prefix = "nexus/cal/dev"
  }
}
