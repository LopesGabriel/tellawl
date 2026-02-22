# ─── Wildcard A record for K8s services ─────────────────────────
# *.k8s.lopesgabriel.dev → K8s instance
# Traefik handles subdomain routing inside the cluster
resource "aws_route53_record" "k8s_wildcard" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "*.${var.k8s_subdomain}.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_instance.k8s.public_ip]
}

# k8s.lopesgabriel.dev → K8s instance (for ArgoCD, Traefik dashboard, etc.)
resource "aws_route53_record" "k8s" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "${var.k8s_subdomain}.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_instance.k8s.public_ip]
}

# infra.lopesgabriel.dev → Infra instance
resource "aws_route53_record" "infra" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "${var.infra_subdomain}.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_instance.infra.public_ip]
}

# grafana.lopesgabriel.dev → Infra instance (Grafana shortcut)
resource "aws_route53_record" "grafana" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = "grafana.${var.domain_name}"
  type    = "A"
  ttl     = 300
  records = [aws_instance.infra.public_ip]
}
