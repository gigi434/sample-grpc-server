#!/bin/bash
set -e

echo "Setting up development environment..."

# Serena MCP サーバーの設定
echo "Configuring Serena MCP server..."
claude mcp add serena -- uvx --from git+https://github.com/oraios/serena serena-mcp-server --context ide-assistant --project /workspaces/${localWorkspaceFolderBasename}

# Python仮想環境の設定
echo "Creating Python virtual environment..."
uv venv
source .venv/bin/activate

# SuperClaudeのインストール
echo "Installing SuperClaude..."
# SuperClaudeパッケージをインストール
uv pip install SuperClaude

# SuperClaudeのセットアップを実行（開発者プロファイルで完全セットアップ）
echo "Setting up SuperClaude framework..."
python3 -m SuperClaude install --profile developer || echo "SuperClaude setup may require manual configuration"

echo "Setup complete!"