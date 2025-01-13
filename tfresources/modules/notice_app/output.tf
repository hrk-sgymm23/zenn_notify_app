output "lambda_arn" {
  value = aws_lambda_function.some_func.arn
}

output "lambda_name" {
  value = aws_lambda_function.some_func.function_name
}