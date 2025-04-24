
PROTOS = \
	common/common.proto \
	fraud_detection/fraud_detection.proto \
	transaction_verification/transaction_verification.proto \
	suggestions/suggestions.proto \
	order_queue/order_queue.proto \
	database/database.proto \
	replica/replica.proto \

# gen:
# 	@for proto in $(PROTOS); do \
# 		echo $$proto; \
# 		protoc --proto_path=pb --go_out=. --go_opt=paths=source_relative \
# 		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
# 		$$proto; \
# 	done 
	
gen_go:
	@for proto in $(PROTOS); do \
		echo $$proto; \
		protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$$proto; \
	done


gen_py:
	@for proto in $(PROTOS); do \
		echo $$proto; \
		python3 -m grpc_tools.protoc -I. --proto_path=. --python_out=. --pyi_out=. --grpc_python_out=. $$proto; \
	done