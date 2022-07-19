proto:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        cmd/heimdall-grpc/proto/token.proto

mockgen:
	mockgen -source=pkg/encryption/aes.go -destination=test/mock_encryption/mock_aes.go
	mockgen -source=pkg/signature/jws.go -destination=test/mock_signature/mock_jws.go