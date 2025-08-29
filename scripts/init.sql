-- データベース初期化スクリプト
-- このスクリプトはDockerコンテナ起動時に自動実行されます

-- データベースが存在しない場合は作成
SELECT 'CREATE DATABASE financial_development'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'financial_development')\gexec

-- 拡張機能を有効化
\c financial_development;

-- UUID生成用の拡張機能
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 暗号化用の拡張機能
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- テキスト検索用の拡張機能（将来の検索機能のため）
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- デフォルトのスキーマに権限を付与
GRANT ALL PRIVILEGES ON DATABASE financial_development TO postgres;
GRANT ALL ON SCHEMA public TO postgres;