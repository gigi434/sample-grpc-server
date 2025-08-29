FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

# Switch to vscode user first
USER vscode

# Install uv as vscode user
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# Add uv to PATH
ENV PATH="/home/vscode/.local/bin:$PATH"

# Create virtual environment and install SuperClaude
RUN /home/vscode/.local/bin/uv venv /home/vscode/.venv && \
    /home/vscode/.local/bin/uv pip install --python /home/vscode/.venv/bin/python SuperClaude

# Add venv to PATH
ENV PATH="/home/vscode/.venv/bin:$PATH"

# Run SuperClaude install (allow failure)
RUN python -m SuperClaude install || true