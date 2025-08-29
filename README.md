# ⚠️⚠️⚠️ 重要な警告 ⚠️⚠️⚠️

## 🤖 このプロジェクトは100% CLAUDE CODEによるペアプログラミングで作成されました 🤖

### ⚡ すべてのコードはCLAUDE CODEとのペアプログラミングによって自動生成されています ⚡

**このリポジトリ内のすべてのコード、設定、ドキュメントは、人間とCLAUDE CODEのペアプログラミングセッションで作成されました。手動でのコーディングは一切行われていません。**

---

# Sample gRPC Server

モジュラーモノリス構造のgRPCサーバー実装

## 🚀 クイックスタート

### 前提条件

- Go 1.20以上
- PostgreSQL 14以上
- Make
- Docker (オプション)

### セットアップ

1. **環境変数の設定**
```bash
cp .env.example .env
# .envファイルを編集して、データベース接続情報を設定
```

2. **依存関係のインストール**
```bash
make deps
```

3. **データベースの起動**

Dockerを使用する場合:
```bash
make db-start
```

既存のPostgreSQLを使用する場合は、`.env`ファイルの接続情報を更新してください。

4. **データベースマイグレーション**
```bash
make migrate
```

5. **シードデータの投入**
```bash
make seed
```

## 📦 プロジェクト構造

```
sample-grpc-server/
├── api/                    # API定義
│   └── proto/              # Protocol Buffers定義
├── cmd/                    # アプリケーションエントリーポイント
│   ├── server/             # gRPCサーバー
│   └── seed/               # データベースシード
├── internal/               # 内部パッケージ
│   ├── config/             # 設定管理
│   ├── modules/            # ビジネスモジュール
│   │   └── user/           # ユーザーモジュール
│   │       ├── domain/     # ドメイン層
│   │       ├── application/# アプリケーション層
│   │       └── infrastructure/ # インフラ層
│   └── shared/             # 共有コンポーネント
├── pkg/                    # 外部公開パッケージ
├── test/                   # テストコード
│   └── fixtures/           # テストデータ
└── Makefile                # ビルド・実行コマンド
```

## 🛠️ Makeコマンド

### 基本コマンド

| コマンド | 説明 |
|---------|------|
| `make help` | 利用可能なコマンドを表示 |
| `make build` | バイナリをビルド |
| `make run` | アプリケーションを実行 |
| `make test` | テストを実行 |
| `make clean` | ビルド成果物を削除 |

### データベース関連

| コマンド | 説明 |
|---------|------|
| `make seed` | シードデータを投入 |
| `make seed-clean` | データをクリアして再シード |
| `make migrate` | マイグレーションを実行 |
| `make db-start` | PostgreSQLコンテナを起動 |
| `make db-stop` | PostgreSQLコンテナを停止 |
| `make db-shell` | PostgreSQLシェルに接続 |

### 開発ツール

| コマンド | 説明 |
|---------|------|
| `make fmt` | コードをフォーマット |
| `make lint` | リンターを実行 |
| `make vet` | go vetを実行 |
| `make check` | すべてのチェックを実行 |
| `make deps` | 依存関係をダウンロード |
| `make dev` | 開発環境を準備 |

### Docker関連

| コマンド | 説明 |
|---------|------|
| `make docker-build` | Dockerイメージをビルド |
| `make docker-up` | Dockerコンテナを起動 |
| `make docker-down` | Dockerコンテナを停止 |
| `make docker-logs` | Dockerログを表示 |

## 👥 シードデータ

`test/fixtures/users.json`に20人のサンプルユーザーが定義されています：

- **管理者ユーザー**: 
  - `admin@example.com` (パスワード: Admin123!@#)
  - `developer@example.com` (パスワード: Dev123!@#)

- **一般ユーザー**: 18人のサンプルユーザー
  - デフォルトパスワード: `Password123!`

## 🔧 設定

設定は環境変数で管理されます。`.env`ファイルで以下の設定が可能：

```env
# データベース設定
DATABASE_DRIVER=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=sample_grpc_server
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=postgres
DATABASE_SSL_MODE=disable
DATABASE_MAX_OPEN_CONNS=25
DATABASE_MAX_IDLE_CONNS=25
DATABASE_CONN_MAX_LIFETIME=300

# サーバー設定
SERVER_PORT=50051
SERVER_HOST=0.0.0.0
```

## 🧪 テスト

```bash
# すべてのテストを実行
make test

# カバレッジレポート付きでテスト
make test-coverage
```

## 📝 開発ワークフロー

1. **新機能の開発**
```bash
make dev          # 開発環境を準備
make run          # サーバーを起動
```

2. **コードの品質チェック**
```bash
make check        # すべてのチェックを実行
```

3. **データベースのリセット**
```bash
make seed-clean   # データをクリアして再シード
```

## 📄 ライセンス

[MIT License](LICENSE)