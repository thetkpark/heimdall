proto:
	protoc --go_out=./cmd/heimdall/proto --go_opt=paths=source_relative \
        --go-grpc_out=./cmd/heimdall/proto --go-grpc_opt=paths=source_relative \
        --proto_path=cmd/heimdall/proto \
        --validate_out="lang=go:." \
        cmd/heimdall/proto/token.proto

mockgen:
	mockgen -source=pkg/encryption/aes.go -destination=test/mock_encryption/mock_aes.go
	mockgen -source=pkg/signature/jws.go -destination=test/mock_signature/mock_jws.go
	mockgen -source=pkg/token/token.go -destination=test/mock_token/mock_token.go

unit-test:
	ginkgo -r

swagger:
	swag init --dir cmd/heimdall --parseDependency --parseInternal