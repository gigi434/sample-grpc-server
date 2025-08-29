# ディレクトリ構成設計書

## 1. 推奨ディレクトリ構成

```
sample-grpc-server/
├── api/                        # API定義
│   └── proto/                  # Protocol Buffers定義
│       ├── common/             # 共通型定義
│       │   ├── pagination.proto
│       │   └── timestamp.proto
│       └── v1/                 # APIバージョン
│           ├── user/           # ユーザーサービス
│           │   └── user.proto
│           └── health/         # ヘルスチェック
│               └── health.proto
│
├── cmd/                        # アプリケーションエントリーポイント
│   └── server/                 # gRPCサーバー
│       └── main.go (main.ts, main.py)
│
├── internal/                   # 内部パッケージ（Go特有、他言語では src/）
│   ├── config/                 # 設定管理
│   │   ├── config.go
│   │   └── validation.go
│   │
│   ├── modules/                # ビジネスモジュール
│   │   ├── user/               # ユーザーモジュール
│   │   │   ├── domain/         # ドメイン層
│   │   │   │   ├── entity/     # エンティティ
│   │   │   │   │   └── user.go
│   │   │   │   ├── repository/ # リポジトリインターフェース
│   │   │   │   │   └── user_repository.go
│   │   │   │   └── service/    # ドメインサービス
│   │   │   │       └── user_service.go
│   │   │   │
│   │   │   ├── application/    # アプリケーション層
│   │   │   │   ├── usecase/    # ユースケース
│   │   │   │   │   ├── create_user.go
│   │   │   │   │   ├── get_user.go
│   │   │   │   │   └── list_users.go
│   │   │   │   └── dto/        # データ転送オブジェクト
│   │   │   │       └── user_dto.go
│   │   │   │
│   │   │   ├── infrastructure/ # インフラストラクチャ層
│   │   │   │   ├── persistence/# 永続化実装
│   │   │   │   │   └── user_repository_impl.go
│   │   │   │   └── grpc/       # gRPCサービス実装
│   │   │   │       └── user_service.go
│   │   │   │
│   │   │   └── module.go       # モジュール定義・DI設定
│   │   │
│   │   └── health/             # ヘルスチェックモジュール
│   │       └── ...             # 同様の構造
│   │
│   ├── shared/                 # 共有コンポーネント
│   │   ├── database/           # データベース関連
│   │   │   ├── connection.go
│   │   │   └── migration.go
│   │   ├── logger/             # ロギング
│   │   │   └── logger.go
│   │   ├── middleware/         # ミドルウェア
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── recovery.go
│   │   ├── errors/             # エラー定義
│   │   │   └── errors.go
│   │   └── utils/              # ユーティリティ
│   │       └── validation.go
│   │
│   └── server/                 # サーバー設定
│       ├── grpc.go             # gRPCサーバー初期化
│       └── interceptor.go      # インターセプター設定
│
├── pkg/                        # 外部公開パッケージ（Go特有）
│   └── generated/              # 生成コード
│       └── api/                # Protocol Buffersから生成
│           └── v1/
│               └── user/
│                   ├── user.pb.go
│                   └── user_grpc.pb.go
│
├── scripts/                    # スクリプト
│   ├── generate.sh             # コード生成
│   ├── migrate.sh              # DB マイグレーション
│   └── test.sh                 # テスト実行
│
├── deployments/                # デプロイメント設定
│   ├── docker/                 # Docker関連
│   │   └── Dockerfile
│   ├── kubernetes/             # Kubernetes マニフェスト
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── docker-compose.yml      # 開発環境用
│
├── test/                       # テストコード
│   ├── integration/            # 統合テスト
│   │   └── user_test.go
│   ├── e2e/                    # E2Eテスト
│   │   └── user_scenario_test.go
│   └── fixtures/               # テストデータ
│       └── users.json
│
├── docs/                       # ドキュメント
│   ├── architecture.md         # アーキテクチャ設計書
│   ├── directory-structure.md  # 本ドキュメント
│   └── api/                    # API ドキュメント
│       └── user.md
│
├── configs/                    # 設定ファイル
│   ├── default.yaml            # デフォルト設定
│   ├── development.yaml        # 開発環境設定
│   ├── staging.yaml            # ステージング環境設定
│   └── production.yaml         # 本番環境設定
│
├── .github/                    # GitHub関連
│   └── workflows/              # GitHub Actions
│       ├── ci.yml
│       └── cd.yml
│
├── .gitignore                  # Git除外設定
├── .env.example                # 環境変数サンプル
├── Makefile                    # ビルド・実行コマンド
├── go.mod                      # Go モジュール定義
├── go.sum                      # Go 依存関係ロック
└── README.md                   # プロジェクト説明

```

## 2. 各ディレクトリの詳細説明

### 2.1 api/proto/
Protocol Buffers定義ファイルを配置します。
- **common/**: 複数のサービスで共有される型定義
- **v1/**: APIバージョンごとに分離
- 各サービスごとにサブディレクトリを作成

### 2.2 cmd/
アプリケーションのエントリーポイントを配置します。
- 各実行可能ファイルごとにサブディレクトリを作成
- main関数のみを含み、ビジネスロジックは含めない

### 2.3 internal/（Go）または src/（他言語）
アプリケーションの主要なコードを配置します。

#### modules/
各ビジネスモジュールを独立したディレクトリとして管理：
- **domain/**: ビジネスロジックとエンティティ
- **application/**: ユースケースとDTO
- **infrastructure/**: 技術的な実装詳細

#### shared/
モジュール間で共有されるコンポーネント：
- データベース接続
- ロギング
- 共通ミドルウェア
- エラーハンドリング

### 2.4 pkg/
外部に公開可能なパッケージ（Go特有）。
生成されたProtocol Buffersコードなどを配置。

### 2.5 scripts/
開発・運用で使用するスクリプト：
- コード生成
- データベースマイグレーション
- テスト実行
- デプロイメント

### 2.6 deployments/
デプロイメント関連の設定：
- Dockerfile
- Kubernetes マニフェスト
- docker-compose.yml（開発環境用）

### 2.7 test/
テストコード：
- **integration/**: 複数コンポーネントの統合テスト
- **e2e/**: エンドツーエンドテスト
- **fixtures/**: テストデータ

### 2.8 configs/
環境別の設定ファイル：
- YAML形式を推奨
- 環境変数でオーバーライド可能な設計

## 3. 言語別の考慮事項

### 3.1 Go言語
```
internal/           # Go特有の内部パッケージ
pkg/                # 公開パッケージ
go.mod, go.sum      # モジュール管理
```

### 3.2 TypeScript/Node.js
```
src/                # internal の代わり
dist/               # ビルド出力
package.json        # 依存関係管理
tsconfig.json       # TypeScript設定
```

### 3.3 Python
```
src/                # internal の代わり
requirements.txt    # 依存関係管理
setup.py            # パッケージ設定
pyproject.toml      # 最新のパッケージ管理
```

## 4. ファイル命名規則

### 4.1 一般的な規則
- **小文字とアンダースコア**: `user_service.go`
- **単数形**: `user.proto`（複数のユーザーを扱う場合でも）
- **明確な名前**: 省略形を避ける

### 4.2 言語別規則
- **Go**: snake_case または camelCase
- **TypeScript**: camelCase、kebab-case（ファイル名）
- **Python**: snake_case

## 5. モジュール間の依存関係

```
┌─────────────────┐
│  Presentation   │
│   (gRPC/REST)   │
└────────┬────────┘
         │
┌────────▼────────┐
│  Application    │
│   (Use Cases)   │
└────────┬────────┘
         │
┌────────▼────────┐
│     Domain      │
│  (Business)     │
└────────┬────────┘
         │
┌────────▼────────┐
│ Infrastructure  │
│  (Technical)    │
└─────────────────┘
```

### 5.1 依存の方向
- 外側から内側への依存のみ許可
- ドメイン層は他の層に依存しない
- インターフェースを通じた依存性逆転

### 5.2 共有コンポーネント
- 横断的関心事は shared/ に配置
- 各モジュールから利用可能
- ビジネスロジックを含めない

## 6. ベストプラクティス

### 6.1 モジュール設計
- **単一責任**: 各モジュールは1つの責任を持つ
- **高凝集・疎結合**: モジュール内は密接に、モジュール間は疎に
- **テスタブル**: 依存性注入により単体テスト可能に

### 6.2 ディレクトリ深さ
- 最大4階層程度に抑える
- 深すぎる階層は複雑性を増す
- 必要に応じてフラット化を検討

### 6.3 生成コード
- pkg/generated/ または build/ に配置
- バージョン管理に含めるかは要検討
- 生成スクリプトを用意

## 7. セキュリティ考慮事項

### 7.1 機密情報
- 環境変数または外部設定サービスを使用
- .env ファイルはサンプルのみコミット
- configs/ には機密情報を含めない

### 7.2 アクセス制御
- internal/ または src/ は外部公開しない
- pkg/ のみ外部利用可能（Go の場合）

## 8. 拡張性の考慮

### 8.1 新規モジュール追加
1. modules/ 配下に新規ディレクトリ作成
2. domain/application/infrastructure の構造を維持
3. 既存モジュールとの依存関係を最小化

### 8.2 バージョニング
- api/proto/v2/ のように新バージョンを追加
- 後方互換性を考慮した設計

## 9. まとめ

この構成により以下が実現されます：
- **明確な責任分離**: 各層・モジュールの役割が明確
- **高い保守性**: 一貫した構造により理解しやすい
- **拡張性**: 新機能追加が容易
- **テスタビリティ**: 各層を独立してテスト可能
- **将来性**: マイクロサービス化への移行が容易

プロジェクトの規模や要件に応じて、この構成をカスタマイズすることを推奨します。