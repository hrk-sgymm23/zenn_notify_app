terraform {
  backend "s3" {
    bucket  = "zenn-notice-app-tfstate-bucket"
    key     = "api/terraform.tfstate"
    region  = "ap-northeast-1"
    encrypt = true
  }
}