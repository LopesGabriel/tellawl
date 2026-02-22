variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "AWS CLI profile"
  type        = string
}

variable "domain_name" {
  description = "Root domain name managed in Route53"
  type        = string
  default     = "lopesgabriel.dev"
}

variable "instance_type" {
  description = "EC2 instance type (must be Graviton/ARM64)"
  type        = string
  default     = "t4g.medium"
}

variable "k8s_subdomain" {
  description = "Subdomain for K8s services (wildcard A record)"
  type        = string
  default     = "k8s"
}

variable "infra_subdomain" {
  description = "Subdomain for the infra instance"
  type        = string
  default     = "infra"
}

variable "my_ip" {
  description = "Your public IP for SSH access (CIDR notation, e.g. 203.0.113.10/32). Use 0.0.0.0/0 to allow from anywhere (not recommended)."
  type        = string
}

variable "root_volume_size" {
  description = "Root EBS volume size in GB"
  type        = number
  default     = 30
}
