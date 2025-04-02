run:
	source .env && go run cmd/server/main.go


		"localhost:50052",
		"localhost:50053",
		"localhost:50054",
1:
	export PORT="50052" && go run cmd/raft/main.go

2:
	export PORT="50053" && go run cmd/raft/main.go 

3:
	export PORT="50054" && go run cmd/raft/main.go 

4:
	export PORT="50055" && go run cmd/raft/main.go 

gen:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pb/order_srv.proto