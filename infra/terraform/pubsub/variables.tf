variable "project_id" {
  type = string
  description = "Your GCP project id."
}

variable "region" {
  type    = string
}

variable "topic_name" {
  type    = string
  description = "Name of the topic to be created for the HelloApi service."
  default = "cymbal-topic"
}

variable "subscription_name" {
  type    = string
  description = "Name of the subscription to be created for the HelloApi service."
  default = "cymbal-subscription"
}

