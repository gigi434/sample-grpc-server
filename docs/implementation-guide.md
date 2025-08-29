# 実装ガイド - ユーザーモジュール

## 1. 概要

このドキュメントでは、モジュラーモノリスgRPCサーバーにおけるユーザーモジュールの実装例を示します。

## 2. Protocol Buffers定義

### api/proto/v1/user/user.proto

```protobuf
syntax = "proto3";

package api.v1.user;

option go_package = "github.com/yourorg/sample-grpc-server/pkg/generated/api/v1/user";

import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";

// ユーザーサービス定義
service UserService {
  // ユーザー作成
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  
  // ユーザー取得
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // ユーザー一覧取得
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  
  // ユーザー更新
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  
  // ユーザー削除
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
}

// ユーザーエンティティ
message User {
  string id = 1;
  string email = 2;
  string name = 3;
  UserRole role = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}

// ユーザーロール
enum UserRole {
  USER_ROLE_UNSPECIFIED = 0;
  USER_ROLE_USER = 1;
  USER_ROLE_ADMIN = 2;
}

// リクエスト/レスポンス定義
message CreateUserRequest {
  string email = 1;
  string name = 2;
  string password = 3;
  UserRole role = 4;
}

message CreateUserResponse {
  User user = 1;
}

message GetUserRequest {
  string id = 1;
}

message GetUserResponse {
  User user = 1;
}

message ListUsersRequest {
  int32 page_size = 1;
  string page_token = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  string next_page_token = 2;
}

message UpdateUserRequest {
  string id = 1;
  string email = 2;
  string name = 3;
  UserRole role = 4;
}

message UpdateUserResponse {
  User user = 1;
}

message DeleteUserRequest {
  string id = 1;
}
```

## 3. ドメイン層実装例

### Go言語の例

#### internal/modules/user/domain/entity/user.go

```go
package entity

import (
    "time"
    "errors"
)

// User ドメインエンティティ
type User struct {
    ID        string
    Email     string
    Name      string
    Password  string // ハッシュ化されたパスワード
    Role      UserRole
    CreatedAt time.Time
    UpdatedAt time.Time
}

// UserRole ユーザーロール
type UserRole int

const (
    UserRoleUnspecified UserRole = iota
    UserRoleUser
    UserRoleAdmin
)

// NewUser ユーザーエンティティの生成
func NewUser(email, name, password string, role UserRole) (*User, error) {
    if err := validateEmail(email); err != nil {
        return nil, err
    }
    
    if err := validatePassword(password); err != nil {
        return nil, err
    }
    
    return &User{
        Email:    email,
        Name:     name,
        Password: password,
        Role:     role,
    }, nil
}

// ビジネスルールの実装
func validateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    // メールアドレスの形式チェック
    return nil
}

func validatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    return nil
}
```

#### internal/modules/user/domain/repository/user_repository.go

```go
package repository

import (
    "context"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/domain/entity"
)

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    List(ctx context.Context, limit int, offset int) ([]*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
}
```

## 4. アプリケーション層実装例

#### internal/modules/user/application/usecase/create_user.go

```go
package usecase

import (
    "context"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/domain/entity"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/domain/repository"
    "github.com/yourorg/sample-grpc-server/internal/shared/logger"
)

// CreateUserUseCase ユーザー作成ユースケース
type CreateUserUseCase struct {
    userRepo repository.UserRepository
    logger   logger.Logger
}

// NewCreateUserUseCase コンストラクタ
func NewCreateUserUseCase(
    userRepo repository.UserRepository,
    logger logger.Logger,
) *CreateUserUseCase {
    return &CreateUserUseCase{
        userRepo: userRepo,
        logger:   logger,
    }
}

// Execute ユースケースの実行
func (uc *CreateUserUseCase) Execute(
    ctx context.Context,
    email, name, password string,
    role entity.UserRole,
) (*entity.User, error) {
    // ビジネスロジック: 既存ユーザーのチェック
    existingUser, _ := uc.userRepo.FindByEmail(ctx, email)
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }
    
    // ユーザーエンティティの生成
    user, err := entity.NewUser(email, name, password, role)
    if err != nil {
        return nil, err
    }
    
    // パスワードのハッシュ化（実際の実装では適切なハッシュ関数を使用）
    user.Password = hashPassword(password)
    
    // リポジトリへの保存
    if err := uc.userRepo.Create(ctx, user); err != nil {
        uc.logger.Error("failed to create user", "error", err)
        return nil, err
    }
    
    uc.logger.Info("user created successfully", "userId", user.ID)
    return user, nil
}
```

## 5. インフラストラクチャ層実装例

#### internal/modules/user/infrastructure/grpc/user_service.go

```go
package grpc

import (
    "context"
    pb "github.com/yourorg/sample-grpc-server/pkg/generated/api/v1/user"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/application/usecase"
)

// UserServiceServer gRPCサービス実装
type UserServiceServer struct {
    pb.UnimplementedUserServiceServer
    createUserUC *usecase.CreateUserUseCase
    getUserUC    *usecase.GetUserUseCase
    listUsersUC  *usecase.ListUsersUseCase
}

// NewUserServiceServer コンストラクタ
func NewUserServiceServer(
    createUserUC *usecase.CreateUserUseCase,
    getUserUC *usecase.GetUserUseCase,
    listUsersUC *usecase.ListUsersUseCase,
) *UserServiceServer {
    return &UserServiceServer{
        createUserUC: createUserUC,
        getUserUC:    getUserUC,
        listUsersUC:  listUsersUC,
    }
}

// CreateUser ユーザー作成
func (s *UserServiceServer) CreateUser(
    ctx context.Context,
    req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
    // ユースケースの実行
    user, err := s.createUserUC.Execute(
        ctx,
        req.GetEmail(),
        req.GetName(),
        req.GetPassword(),
        convertRole(req.GetRole()),
    )
    if err != nil {
        return nil, err
    }
    
    // レスポンスの生成
    return &pb.CreateUserResponse{
        User: convertUserToProto(user),
    }, nil
}

// GetUser ユーザー取得
func (s *UserServiceServer) GetUser(
    ctx context.Context,
    req *pb.GetUserRequest,
) (*pb.GetUserResponse, error) {
    user, err := s.getUserUC.Execute(ctx, req.GetId())
    if err != nil {
        return nil, err
    }
    
    return &pb.GetUserResponse{
        User: convertUserToProto(user),
    }, nil
}
```

## 6. 依存性注入の設定

#### internal/modules/user/module.go

```go
package user

import (
    "github.com/yourorg/sample-grpc-server/internal/modules/user/application/usecase"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/infrastructure/grpc"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/infrastructure/persistence"
    "github.com/yourorg/sample-grpc-server/internal/shared/database"
    "github.com/yourorg/sample-grpc-server/internal/shared/logger"
)

// Module ユーザーモジュール
type Module struct {
    GRPCService *grpc.UserServiceServer
}

// NewModule モジュールの初期化
func NewModule(db *database.DB, logger logger.Logger) *Module {
    // リポジトリの初期化
    userRepo := persistence.NewUserRepository(db)
    
    // ユースケースの初期化
    createUserUC := usecase.NewCreateUserUseCase(userRepo, logger)
    getUserUC := usecase.NewGetUserUseCase(userRepo, logger)
    listUsersUC := usecase.NewListUsersUseCase(userRepo, logger)
    
    // gRPCサービスの初期化
    grpcService := grpc.NewUserServiceServer(
        createUserUC,
        getUserUC,
        listUsersUC,
    )
    
    return &Module{
        GRPCService: grpcService,
    }
}
```

## 7. サーバー起動

#### cmd/server/main.go

```go
package main

import (
    "log"
    "net"
    
    "google.golang.org/grpc"
    pb "github.com/yourorg/sample-grpc-server/pkg/generated/api/v1/user"
    "github.com/yourorg/sample-grpc-server/internal/config"
    "github.com/yourorg/sample-grpc-server/internal/modules/user"
    "github.com/yourorg/sample-grpc-server/internal/shared/database"
    "github.com/yourorg/sample-grpc-server/internal/shared/logger"
)

func main() {
    // 設定の読み込み
    cfg := config.Load()
    
    // ロガーの初期化
    logger := logger.New(cfg.LogLevel)
    
    // データベース接続
    db, err := database.Connect(cfg.Database)
    if err != nil {
        log.Fatal("failed to connect to database:", err)
    }
    defer db.Close()
    
    // モジュールの初期化
    userModule := user.NewModule(db, logger)
    
    // gRPCサーバーの起動
    lis, err := net.Listen("tcp", cfg.Server.Port)
    if err != nil {
        log.Fatal("failed to listen:", err)
    }
    
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, userModule.GRPCService)
    
    logger.Info("gRPC server starting", "port", cfg.Server.Port)
    if err := s.Serve(lis); err != nil {
        log.Fatal("failed to serve:", err)
    }
}
```

## 8. Makefile例

```makefile
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make proto    - Generate code from proto files"
	@echo "  make build    - Build the server"
	@echo "  make run      - Run the server"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean generated files"

.PHONY: proto
proto:
	@echo "Generating proto files..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/v1/**/*.proto

.PHONY: build
build:
	@echo "Building server..."
	go build -o bin/server cmd/server/main.go

.PHONY: run
run:
	@echo "Running server..."
	go run cmd/server/main.go

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf bin/ pkg/generated/
```

## 9. テスト実装例

#### internal/modules/user/application/usecase/create_user_test.go

```go
package usecase_test

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/application/usecase"
    "github.com/yourorg/sample-grpc-server/internal/modules/user/domain/entity"
)

// モックリポジトリ
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func TestCreateUserUseCase_Execute(t *testing.T) {
    ctx := context.Background()
    
    t.Run("successful user creation", func(t *testing.T) {
        // モックの設定
        mockRepo := new(MockUserRepository)
        mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, nil)
        mockRepo.On("Create", ctx, mock.Anything).Return(nil)
        
        // ユースケースの実行
        uc := usecase.NewCreateUserUseCase(mockRepo, nil)
        user, err := uc.Execute(ctx, "test@example.com", "Test User", "password123", entity.UserRoleUser)
        
        // アサーション
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "test@example.com", user.Email)
        mockRepo.AssertExpectations(t)
    })
    
    t.Run("user already exists", func(t *testing.T) {
        // 既存ユーザーが存在する場合
        existingUser := &entity.User{Email: "test@example.com"}
        mockRepo := new(MockUserRepository)
        mockRepo.On("FindByEmail", ctx, "test@example.com").Return(existingUser, nil)
        
        // ユースケースの実行
        uc := usecase.NewCreateUserUseCase(mockRepo, nil)
        user, err := uc.Execute(ctx, "test@example.com", "Test User", "password123", entity.UserRoleUser)
        
        // アサーション
        assert.Error(t, err)
        assert.Nil(t, user)
        mockRepo.AssertExpectations(t)
    })
}
```

## 10. 設定ファイル例

#### configs/default.yaml

```yaml
server:
  port: ":50051"
  timeout: 30s

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  name: "grpc_server"
  user: "postgres"
  password: "password"
  sslmode: "disable"
  max_open_conns: 25
  max_idle_conns: 25
  conn_max_lifetime: 5m

logger:
  level: "info"
  format: "json"

auth:
  jwt_secret: "your-secret-key"
  token_expiry: 24h

monitoring:
  metrics_enabled: true
  metrics_port: ":9090"
  tracing_enabled: true
  tracing_endpoint: "http://localhost:14268/api/traces"
```

## 11. まとめ

このガイドでは、モジュラーモノリスgRPCサーバーにおけるユーザーモジュールの実装例を示しました。

### 主なポイント

1. **レイヤー分離**: ドメイン、アプリケーション、インフラストラクチャの明確な分離
2. **依存性注入**: インターフェースを通じた疎結合な設計
3. **テスタビリティ**: モックを使用した単体テストの実装
4. **設定管理**: 環境別の設定ファイル
5. **エラーハンドリング**: 適切なエラーの伝搬と処理

この構成により、保守性が高く、拡張可能なgRPCサーバーを構築できます。