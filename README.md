# Heimdall

[![Docker](https://github.com/thetkpark/heimdall/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/thetkpark/heimdall/actions/workflows/docker-publish.yml) ![Test](https://github.com/thetkpark/heimdall/actions/workflows/unit-test.yml/badge.svg) ![GitHub](https://img.shields.io/github/license/thetkpark/heimdall) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thetkpark/heimdall)

> **Heimdall** is a [god](https://en.wikipedia.org/wiki/Æsir) who keeps watch for invaders and the onset of [Ragnarök](https://en.wikipedia.org/wiki/Ragnarök) from his dwelling [Himinbjörg](https://en.wikipedia.org/wiki/Himinbjörg), where the burning rainbow bridge [Bifröst](https://en.wikipedia.org/wiki/Bifröst)meets the sky.

## Features

- Generate the Json Web Signature of the payload
- Encrypt the payload before signing it for confidentiality
- Verify and parse the payload from the given token
- Verify and set the payload data to HTTP response headers to be used as authentication service
- Token authentication and generation via REST API
- Token generation via gRPC

## Usage

### Environment Variable

| Name                   | Required? | Default value | Note                                      |
| ---------------------- | --------- | ------------- | ----------------------------------------- |
| JWS_SECRET_KEY         | YES       |               |                                           |
| PAYLOAD_ENCRYPTION_KEY |           |               | If omitted, payload will not be encrypted |
| TOKEN_VALID_TIME       |           |               |                                           |
| SENTRY_DSN             |           |               |                                           |
| MODE                   |           | development   |                                           |
| GIN_MODE               |           | debug         |                                           |
| GIN_PORT               |           | 8080          |                                           |
| GRPC_PORT              |           | 5050          |                                           |

### Docker

```shell
docker run -p 8080:8080 -p 5050:5050 -e JWS_SECRET_KEY=SecretKey thetkpark/heimdall
```

### API Specification

#### REST API

> Swagger UI is availble at `/swagger/index.html`

#### gRPC

> Please look at the Protocol Buffers file in `cmd/heimdall/proto/token.proto`
