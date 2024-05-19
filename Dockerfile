# MySQLのベースイメージを使用
FROM mysql:8.0 AS mysql

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

RUN apt-get update && apt-get install -y tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" > /etc/timezone

# 初期化スクリプトをコピー
COPY ./init.sql /docker-entrypoint-initdb.d/

# Goビルド環境をベースイメージとして使用
FROM golang:1.21 AS builder

# ワークディレクトリを設定
WORKDIR /app

# go.mod と go.sum ファイルをコピー
COPY go.mod .
COPY go.sum .

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY ./cmd/mynote ./cmd/mynote

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mynote ./cmd/mynote/main.go

# 新しいステージを開始し、Alpine Linuxイメージをベースとする
FROM alpine:latest

# タイムゾーンを設定
ENV TZ=Asia/Tokyo
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" > /etc/timezone

WORKDIR /root/

# ビルドした実行ファイルをAlpineコンテナにコピー
COPY --from=builder /app/mynote .

COPY ./internal/views/template ./template

# コンテナが起動したときに実行されるコマンド
CMD ["./mynote"]
