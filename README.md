# Ollama CLI

A simple command-line chat client for Ollama models.

## Requirements

- Go 1.25.1+
- Ollama server running locally (`http://localhost:11434`)

## Installation

```bash
git clone https://github.com/teasec4/ollama-go-cli.git
cd ollama-cli
go mod download
```

## Usage

```bash
go run ./cmd/chat
```

## Configuration

Environment variables:
- `OLLAMA_URL` - Ollama API endpoint (default: `http://localhost:11434`)
- `OLLAMA_MODEL` - Model name (default: `llama3:latest`)

## Controls

- Enter - Send message
- PgUp/↑ - Scroll up
- PgDn/↓ - Scroll down
- Ctrl+C - Exit

## License

MIT
