output "service_url" {
  description = "Cloud Run service URL"
  value       = module.cloud_run.service_url
}

output "service_name" {
  description = "Cloud Run service name"
  value       = module.cloud_run.service_name
}

output "data_bucket" {
  description = "GCS bucket for SQLite data"
  value       = google_storage_bucket.cal_data.name
}

output "subscribe_url" {
  description = "Calendar subscription URL template"
  value       = "webcal://${trimprefix(module.cloud_run.service_url, "https://")}/cal/{token}.ics"
}
