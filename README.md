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
