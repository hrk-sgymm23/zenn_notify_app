# ベースイメージとして公式のGolangイメージを使用
FROM golang:1.20-alpine

# 作業ディレクトリを設定
WORKDIR /app

# Go Modulesの使用を許可（必要な場合）
# ENV GO111MODULE=on

# ローカルのモジュールキャッシュを最適化（必要な場合）
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# go.mod と go.sum を先にコピーして依存関係を解決
COPY src/go.mod ./
RUN go mod download

# ソースコードをコンテナにコピー
COPY src/ .

# アプリケーションをビルド
RUN go build -o zenn_notify_app

# 実行可能ファイルをデフォルトのコマンドとして設定
CMD ["./zenn_notify_app"]