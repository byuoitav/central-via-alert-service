terraform {
  backend "s3" {
    bucket         = "terraform-state-storage-586877430255"
    dynamodb_table = "terraform-state-lock-586877430255"
    region         = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "central-via-alert-service.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
  config_path = "~/.kube/config"
}

// pull all env vars out of ssm
data "aws_ssm_parameter" "prd_couch_address" {
  name = "/env/couch-address"
}

data "aws_ssm_parameter" "prd_couch_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "prd_couch_password" {
  name = "/env/couch-password"
}

resource "aws_s3_bucket" "bucket" {
  bucket = "av-flight-deck"

  tags = {
    repo_url         = "https://github.com/byuoitav/central-via-alert-service"
    team             = "AV Engineering"
    data-sensitivity = "confidential"
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

data "aws_iam_policy_document" "policy" {
  statement {
    actions = [
      "s3:ListAllMyBuckets",
      "s3:GetBucketLocation",
    ]
    resources = [
      "arn:aws:s3:::*"
    ]
  }

  statement {
    actions = [
      "s3:ListBucket",
    ]
    resources = [
      "arn:aws:s3:::${aws_s3_bucket.bucket.id}"
    ]
  }

  statement {
    actions = [
      "s3:*",
    ]

    resources = [
      "arn:aws:s3:::${aws_s3_bucket.bucket.id}",
      "arn:aws:s3:::${aws_s3_bucket.bucket.id}/*",
    ]
  }
}

module "prd_deployment" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "central-via-alert-service"
  image          = "byuoitav/central-via-alert-service"
  image_version  = "latest"
  container_port = 8040
  repo_url       = "https://github.com/byuoitav/central-via-alert-service"

  // optional
  iam_policy_doc = data.aws_iam_policy_document.policy.json
  replicas       = 3
  public_urls    = ["via-alert.av.byu.edu"]
  container_env = {
    "AWS_BUCKET_REGION"          = aws_s3_bucket.bucket.region
    "AWS_DEPLOYMENT_KEY"         = data.aws_ssm_parameter.deployment_key.value
    "DB_ADDRESS"                 = data.aws_ssm_parameter.prd_couch_address.value
    "DB_USERNAME"                = data.aws_ssm_parameter.prd_couch_username.value
    "DB_PASSWORD"                = data.aws_ssm_parameter.prd_couch_password.value
    "DOCKER_GITHUB_PASSWORD"     = data.aws_ssm_parameter.docker_github_password.value
    "DOCKER_GITHUB_USERNAME"     = data.aws_ssm_parameter.docker_github_username.value
    "ELASTIC_API_EVENTS"         = data.aws_ssm_parameter.elk_event_api.value
    "LDAP_PASSWORD"              = data.aws_ssm_parameter.ldap_password.value
    "LDAP_USERNAME"              = data.aws_ssm_parameter.ldap_username.value
    "PI_SSH_USERNAME"            = data.aws_ssm_parameter.pi_username.value
    "RASPI_DEPLOYMENT_S3_BUCKET" = aws_s3_bucket.bucket.id
    "STOP_REPLICATION"           = "true"
  }
}
