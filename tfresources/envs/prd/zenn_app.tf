# IAMリソース関連

# Lambda実行ロール
resource "aws_iam_role" "lambda_role" {
  name               = "${var.common_name}-lambda-role"
  assume_role_policy = file("${path.module}/policies/lambda_assume_policy.json")
}

resource "aws_iam_policy" "lambda_policy" {
  name   = "${var.common_name}-lambda-policy"
  policy = file("${path.module}/policies/lambda_policy.json")
}

resource "aws_iam_role_policy_attachment" "name" {
  role       = aws_iam_role.lambda_role.id
  policy_arn = aws_iam_policy.lambda_policy.arn
}

# Scheduler用ロール
resource "aws_iam_role" "scheduler_assume_role" {
  name               = "${var.common_name}-lambda-scheduler-assume-role"
  assume_role_policy = file("${path.module}/policies/eb_assume_policy.json")
}


# lambdaリソース関連
resource "terraform_data" "build_lambda" {
  // ファイルに変更があった場合のみ以下を実行
  triggers_replace = {
    file_content = md5(file("../../../src/main.go")) // 相対パスでLambdaファイルを指定
    // ファイルでなくディレクトリを指定したい場合は以下
    // file_content = sha1(join("", [for f in fileset("../lambda", "*") : filesha1("../lambda/${f}")])) 
  }

  provisioner "local-exec" {
    command = "cd ../../../src/ && GOOS=linux GOARCH=amd64 go build -o ./bootstrap" // ファイルの出力場所を指定する
  }
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "../../../src/bootstrap"
  output_path = "${path.module}/go.zip"
  // path.moduleはこのモジュールが配置されているディレクトリ

  depends_on = [terraform_data.build_lambda]
}

module "zenn_app" {
  source        = "../../modules/notice_app"
  file_name     = "${path.module}/go.zip"
  function_name = "zenn_app"
  handler       = "bootstrap"
  iam_role_arn  = aws_iam_role.lambda_role.arn
  common_name   = "zenn_app"
  code_hash     = data.archive_file.lambda.output_base64sha256
  environment   = "prd"
}

resource "aws_scheduler_schedule" "exec_lambda" {
  name                = "${var.common_name}-exec-${var.environment}"
  schedule_expression = "cron(0 23 * * ? *)"
  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn      = module.zenn_app.lambda_arn
    role_arn = aws_iam_role.scheduler_assume_role.arn
  }
}