# マルチステージビルド - ビルダーステージ
FROM golang:1.21-alpine AS builder

# ビルドに必要なパッケージをインストール
RUN apk add --no-cache git make gcc musl-dev

# 作業ディレクトリを設定
WORKDIR /build

# 依存関係のキャッシュのためにgo.modとgo.sumを先にコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Protocol Buffersのコード生成（必要な場合）
# RUN make proto

# バイナリをビルド
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o server cmd/server/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o seed cmd/seed/main.go

# マルチステージビルド - 実行ステージ
FROM alpine:3.19

# 実行に必要な最小限のパッケージをインストール
RUN apk add --no-cache ca-certificates tzdata postgresql-client

# タイムゾーンを設定
ENV TZ=Asia/Tokyo

# 非rootユーザーを作成
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 作業ディレクトリを設定
WORKDIR /app

# ビルダーステージから実行可能ファイルをコピー
COPY --from=builder /build/server /app/server
COPY --from=builder /build/seed /app/seed

# 設定ファイルとフィクスチャをコピー
COPY --from=builder /build/.env.example /app/.env.example
COPY --from=builder /build/test/fixtures /app/test/fixtures

# 実行可能ファイルに権限を設定
RUN chmod +x /app/server /app/seed

# 所有権を変更
RUN chown -R appuser:appuser /app

# 非rootユーザーに切り替え
USER appuser

# ヘルスチェック
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD grpc_health_probe -addr=:50051 || exit 1

# gRPCポートを公開
EXPOSE 50051

# サーバーを起動
CMD ["/app/server"]