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

Or build and run:

```bash
go build ./cmd/chat -o ollama-chat
./ollama-chat
```

## Features

- 💬 Interactive chat with Ollama models
- 📊 Token usage tracking
- 🎨 Colored console output
- 📝 Session management

## Configuration

Edit `internal/config/config.go`:

```go
OllamaURL:     "http://localhost:11434"  // Ollama API endpoint
Model:         "gpt-oss:20b"              // Model name
AssistantName: "Max"                      // Assistant name
```

## Project Structure

```
.
├── cmd/
│   └── chat/
│       └── main.go              # Entry point
├── internal/
│   ├── client/
│   │   └── ollama.go            # Ollama API client
│   ├── chat/
│   │   └── session.go           # Chat session management
│   ├── config/
│   │   └── config.go            # Configuration
│   └── ui/
│       ├── console.go           # Console rendering
│       └── interactive.go       # Interactive loop
├── go.mod
├── go.sum
└── README.md
```

## Commands

- Type your message and press Enter to send
- Type `exit` or `quit` to exit
- Token count is displayed in the top-right corner

## License

MIT
