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
  host        = data.aws_ssm_parameter.eks_cluster_endpoint.value
  config_path = "~/.kube/config"
}

// pull all env vars out of ssm
data "aws_ssm_parameter" "prd_couch_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "prd_couch_password" {
  name = "/env/couch-password"
}

data "aws_ssm_parameter" "opa_url" {
  name = "/env/opa-url"
}

data "aws_ssm_parameter" "opa_token" {
  name = "/env/viaalert/opa-token"
}

module "prd_deployment" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "central-via-alert-service"
  image          = "docker.pkg.github.com/byuoitav/central-via-alert-service/central-via-alert-service-dev"
  image_version  = "d22d8d3"
  container_port = 8040
  repo_url       = "https://github.com/byuoitav/central-via-alert-service"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["via-alert.av.byu.edu"]
  container_args = [
    "--username", data.aws_ssm_parameter.prd_couch_username.value,
    "--password", data.aws_ssm_parameter.prd_couch_password.value,
    "-a", data.aws_ssm_parameter.opa_url.value,
    "-t", data.aws_ssm_parameter.opa_token.value,
    "--port", "8040",
    "-L", "-1",
  ]
}
