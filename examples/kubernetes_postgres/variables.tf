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

variable "IONOS_REGISTRY_TOKEN" {
  type      = string
  sensitive = true

  validation {
    condition     = length(var.IONOS_REGISTRY_TOKEN) > 0
    error_message = "Variable IONOS_REGISTRY_TOKEN must not be empty."
  }
}

variable "KUBERNETES_VERSION" {
  type    = string
  default = "1.31.2"
}
