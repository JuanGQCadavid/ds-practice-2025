PROTOS = \
	utils/pb/fraud_detection/fraud_detection.proto \
	utils/pb/transaction_verification/transaction_verification.proto \
	utils/pb/suggestions/suggestions.proto

SPECIFIC = \
	utils/pb/fraud_detection/fraud_detection.proto

run:
	docker compose --env-file .env up --build

# python3 -m grpc_tools.protoc -I. --python_out=. --pyi_out=. --grpc_python_out=. $$proto;

gen_go:
	@for proto in $(SPECIFIC); do \
		echo $$proto; \
		protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$$proto; \
	done