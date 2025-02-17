PROTOS = \
	utils/pb/fraud_detection/fraud_detection.proto \
	utils/pb/transaction_verification/transaction_verification.proto

run:
	docker compose up --build

gen:
	@for proto in $(PROTOS); do \
		# echo $$proto; \
		# protoc --go_out=. --go_opt=paths=source_relative \
		# --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		# $$proto; \
		python3 -m grpc_tools.protoc -I. --python_out=. --pyi_out=. --grpc_python_out=. $$proto; \
	done