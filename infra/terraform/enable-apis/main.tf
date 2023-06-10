module "project-services" {
  source  = "terraform-google-modules/project-factory/google//modules/project_services"
  # version = 
  project_id                  = var.project_id
  activate_apis = [
    "compute.googleapis.com",
    "iam.googleapis.com",
    "container.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "servicemanagement.googleapis.com",
    "artifactregistry.googleapis.com",
    "sourcerepo.googleapis.com",
    "clouddeploy.googleapis.com",
    "containersecurity.googleapis.com",
    "pubsub.googleapis.com",
  ]
}