# ─── Security Group: K8s Instance ───────────────────────────────
resource "aws_security_group" "k8s" {
  name        = "tellawl-k8s-sg"
  description = "Security group for K8s (k3s + Traefik) instance"
  vpc_id      = data.aws_vpc.default.id

  # SSH
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  # HTTP (Traefik)
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTPS (Traefik)
  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # K8s API
  ingress {
    description = "K8s API (kubectl)"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "tellawl-k8s-sg"
  }
}

# ─── Security Group: Infra Instance ────────────────────────────
resource "aws_security_group" "infra" {
  name        = "tellawl-infra-sg"
  description = "Security group for infra instance (Postgres, Kafka, Grafana)"
  vpc_id      = data.aws_vpc.default.id

  # SSH
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  # PostgreSQL — only from K8s SG and your IP
  ingress {
    description     = "PostgreSQL from K8s"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.k8s.id]
  }

  ingress {
    description = "PostgreSQL from admin"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  # Kafka — only from K8s SG and your IP
  ingress {
    description     = "Kafka from K8s"
    from_port       = 9092
    to_port         = 9092
    protocol        = "tcp"
    security_groups = [aws_security_group.k8s.id]
  }

  ingress {
    description = "Kafka from admin"
    from_port   = 9092
    to_port     = 9092
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  ingress {
    description = "Redpanda Console from admin"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  # Grafana (LGTM stack)
  ingress {
    description = "Grafana"
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = [var.my_ip]
  }

  # OTLP gRPC — from K8s SG
  ingress {
    description     = "OTLP gRPC from K8s"
    from_port       = 4317
    to_port         = 4317
    protocol        = "tcp"
    security_groups = [aws_security_group.k8s.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "tellawl-infra-sg"
  }
}

# ─── EC2 Instance: K8s (k3s + Traefik) ─────────────────────────
resource "aws_instance" "k8s" {
  ami                    = data.aws_ami.ubuntu_arm64.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.tellawl.key_name
  vpc_security_group_ids = [aws_security_group.k8s.id]
  subnet_id              = data.aws_subnets.default.ids[2]

  associate_public_ip_address = true

  root_block_device {
    volume_size = var.root_volume_size
    volume_type = "gp3"
  }

  user_data = file("${path.module}/scripts/k8s-init.sh")

  tags = {
    Name = "tellawl-k8s"
    Role = "kubernetes"
  }
}

# ─── EC2 Instance: Infra (Postgres, Kafka, Grafana) ────────────
resource "aws_instance" "infra" {
  ami                    = data.aws_ami.ubuntu_arm64.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.tellawl.key_name
  vpc_security_group_ids = [aws_security_group.infra.id]
  subnet_id              = data.aws_subnets.default.ids[2]

  associate_public_ip_address = true

  root_block_device {
    volume_size = var.root_volume_size
    volume_type = "gp3"
  }

  user_data = file("${path.module}/scripts/infra-init.sh")

  tags = {
    Name = "tellawl-infra"
    Role = "infrastructure"
  }
}
