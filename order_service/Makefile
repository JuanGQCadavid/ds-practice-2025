run:
	source .env && go run cmd/server/main.go

1:
	source .env \
	&& export PORT="50501"  \
	&& export order_service_1_dns="localhost:50502"  \
	&& export order_service_2_dns="localhost:50503"  \
	&& go run cmd/raft/main.go

2:
	source .env \
	&& export PORT="50502"  \
	&& export order_service_1_dns="localhost:50501"  \
	&& export order_service_2_dns="localhost:50503"  \
	&& go run cmd/raft/main.go

3:
	source .env \
	&& export PORT="50503"  \
	&& export order_service_1_dns="localhost:50501"  \
	&& export order_service_2_dns="localhost:50502"  \
	&& go run cmd/raft/main.go


gen:
	protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pb/order_srv.proto