variable "IONOS_S3_ACCESS_KEY" {
  type      = string
  sensitive = true

  validation {
    condition     = length(var.IONOS_S3_ACCESS_KEY) > 0
    error_message = "Variable IONOS_S3_ACCESS_KEY must not be empty."
  }
}

variable "IONOS_S3_SECRET_KEY" {
  type      = string
  sensitive = true

  validation {
    condition     = length(var.IONOS_S3_SECRET_KEY) > 0
    error_message = "Variable IONOS_S3_SECRET_KEY must not be empty."
  }
}
