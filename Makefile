proto:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        cmd/heimdall/proto/token.proto

mockgen:
	mockgen -source=pkg/encryption/aes.go -destination=test/mock_encryption/mock_aes.go
	mockgen -source=pkg/signature/jws.go -destination=test/mock_signature/mock_jws.go
	mockgen -source=pkg/token/token.go -destination=test/mock_token/mock_token.go

unit-test:
	ginkgo -r