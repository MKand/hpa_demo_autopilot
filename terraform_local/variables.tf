variable "project_id" {
  type = string
}

variable "region" {
  type    = string
}

variable "zone" {
  type    = string
}

variable "topic_name" {
  type    = string
  default = "cymbal-topic"
}

variable "subscription_name" {
  type    = string
  default = "cymbal-subscription"
}