
PROTOS = \
	pb/fraud_detection/fraud_detection.proto \
	pb/transaction_verification/transaction_verification.proto

gen:
	@for proto in $(PROTOS); do \
		echo $$proto; \
		protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$$proto; \
	done