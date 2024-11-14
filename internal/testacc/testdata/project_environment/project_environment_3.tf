resource "infisical_project" "test_project" {
  name = var.project_name
  slug = var.project_slug
}

variable "project_name" {
  type = string
}

variable "project_slug" {
  type = string
}

variable "environment_slug" {
  type = string
}

resource "infisical_project_environment" "uat-2" {
  name       = "${var.environment_slug}-2"
  project_id = infisical_project.test_project.id
  slug       = "${var.environment_slug}-2"
  position   = 1
}
