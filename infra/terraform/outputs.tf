output "k8s_public_ip" {
  description = "Public IP of the K8s instance"
  value       = aws_instance.k8s.public_ip
}

output "infra_public_ip" {
  description = "Public IP of the infra instance"
  value       = aws_instance.infra.public_ip
}

output "k8s_dns" {
  description = "DNS name for the K8s instance"
  value       = "${var.k8s_subdomain}.${var.domain_name}"
}

output "k8s_wildcard_dns" {
  description = "Wildcard DNS for K8s services"
  value       = "*.${var.k8s_subdomain}.${var.domain_name}"
}

output "infra_dns" {
  description = "DNS name for the infra instance"
  value       = "${var.infra_subdomain}.${var.domain_name}"
}

output "grafana_dns" {
  description = "DNS name for Grafana"
  value       = "grafana.${var.domain_name}"
}

output "ssh_key_file" {
  description = "Path to the generated SSH private key"
  value       = local_file.ssh_private_key.filename
}

output "ssh_k8s" {
  description = "SSH command for K8s instance"
  value       = "ssh -i ${local_file.ssh_private_key.filename} ubuntu@${aws_instance.k8s.public_ip}"
}

output "ssh_infra" {
  description = "SSH command for infra instance"
  value       = "ssh -i ${local_file.ssh_private_key.filename} ubuntu@${aws_instance.infra.public_ip}"
}

output "kubeconfig_command" {
  description = "Command to copy kubeconfig from K8s instance"
  value       = "scp -i ${local_file.ssh_private_key.filename} ubuntu@${aws_instance.k8s.public_ip}:/etc/rancher/k3s/k3s.yaml ./kubeconfig.yaml"
}
