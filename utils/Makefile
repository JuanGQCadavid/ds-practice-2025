
PROTOS = \
	pb/fraud_detection/fraud_detection.proto \
	pb/transaction_verification/transaction_verification.proto

gen-fraud:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pb/fraud_detection/fraud_detection.proto


# gen_all:
#     @- $(foreach proto_path,$(PROTOS), \
# 		$(protoc --go_out=. --go_opt=paths=source_relative \
# 			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
# 			$(proto_path))
#     )


gen:
	@for proto_path in $(PROTOS); do \
		echo $$proto_path
	done



dude:
	@for proto in $(PROTOS); do \
		echo $$proto; \
		protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$$proto; \
	done

gen_all:
	@for proto_path in $(PROTOS); do \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			$$proto_path; \
	done