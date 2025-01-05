resource "aws_s3_bucket" "example" {
  bucket = "my-tf-testhoge-bucket"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}