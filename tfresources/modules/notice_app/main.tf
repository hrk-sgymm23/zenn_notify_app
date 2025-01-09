resource "aws_lambda_function" "some_func" {
  filename      = var.file_name
  function_name = var.function_name
  role          = var.iam_role_arn
  handler       = var.handler
  memory_size   = 128
  timeout       = 300

  source_code_hash = var.code_hash

  runtime = "provided.al2023"

  // build_lambdaしないと失敗するので、depends_onでbuild_lambdaがあれば実行されるようにする
#   depends_on = [terraform_data.build_lambda]
  environment {
    variables = var.environments_variables
  }
}