variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "dea-noctua"
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "cloud_run_image" {
  description = "Cloud Run container image URL"
  type        = string
  default     = "gcr.io/dea-noctua/nexus-cal-dev:latest"
}

variable "service_account_email" {
  description = "Existing service account email (shared with portal)"
  type        = string
}
