run:
	source .env && go run cmd/server/main.go

gen:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pb/order_srv.proto