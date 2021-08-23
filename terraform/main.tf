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
data "aws_ssm_parameter" "prd_couch_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "prd_couch_password" {
  name = "/env/couch-password"
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
  public_urls    = ["via-alert.av.byu.edu"]
  container_args = [
    "--username", data.aws_ssm_parameter.prd_couch_username.value,
    "--password", data.aws_ssm_parameter.prd_couch_password.value,
    "--port", "8040", 
  ]
}
