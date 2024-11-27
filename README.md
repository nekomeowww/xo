# xo

[![Go Reference](https://pkg.go.dev/badge/github.com/nekomeowww/xo.svg)](https://pkg.go.dev/github.com/nekomeowww/xo)
![](https://github.com/nekomeowww/xo/actions/workflows/ci.yml/badge.svg)
[![](https://goreportcard.com/badge/github.com/nekomeowww/xo)](https://goreportcard.com/report/github.com/nekomeowww/xo)

🪐 Universal external Golang utilities, implementations, and even experimental coding patterns, designs

## Development

Clone the repository:

```shell
git clone https://github.com/nekomeowww/xo
cd xo
```

Prepare the dependencies:

```shell
go mod tidy
```

> [!NOTE]
> If you want to work with `protobufs/testpb` directory and generate new Golang code, you need to install the [`buf`](https://buf.build/docs/installation) tool.
>
> ```shell
> cd protobufs
> buf dep update
> buf generate --path ./testpb
> ```

## 🤠 Spec

GoDoc: [https://godoc.org/github.com/nekomeowww/xo](https://godoc.org/github.com/nekomeowww/xo)

## 👪 Other family members of `anyo`

- [nekomeowww/fo](https://github.com/nekomeowww/fo): Functional programming utility library for Go
- [nekomeowww/bo](https://github.com/nekomeowww/bo): BootKit for easily bootstrapping multi-goroutine applications, CLIs
- [nekomeowww/tgo](https://github.com/nekomeowww/tgo): Telegram bot framework for Go
- [nekomeowww/wso](https://github.com/nekomeowww/wso): WebSocket utility library for Go

## 🎆 Other cool related Golang projects I made & maintained

- [nekomeowww/factorio-rcon-api](https://github.com/nekomeowww/factorio-rcon-api): Fully implemented wrapper for Factorio RCON as API
- [Kollama - Ollama Operator](https://github.com/knoway-dev/knoway): Kubernetes Operator for managing Ollama instances across multiple clusters
- [lingticio/llmg](https://github.com/lingticio/llmg): LLM Gateway with gRPC, WebSocket, and RESTful API adapters included.
- [nekomeowww/hyphen](https://github.com/nekomeowww/hyphen): An elegant URL Shortener service
- [nekomeowww/insights-bot](https://github.com/nekomeowww/insights-bot): Webpage summary & chat history recap bot for Telegram
