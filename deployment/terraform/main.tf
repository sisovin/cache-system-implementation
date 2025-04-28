provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "cache_node" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"

  tags = {
    Name = "CacheNode"
  }
}

resource "aws_elasticache_cluster" "cache_cluster" {
  cluster_id           = "cache-cluster"
  engine               = "redis"
  node_type            = "cache.t2.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"

  tags = {
    Name = "CacheCluster"
  }
}

resource "aws_security_group" "cache_sg" {
  name        = "cache_sg"
  description = "Allow inbound traffic to cache nodes"

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_elasticache_subnet_group" "cache_subnet_group" {
  name       = "cache_subnet_group"
  subnet_ids = ["subnet-12345678", "subnet-87654321"]

  tags = {
    Name = "CacheSubnetGroup"
  }
}
