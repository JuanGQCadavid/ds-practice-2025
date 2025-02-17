from datetime import datetime
from concurrent import futures
import grpc
import json
import sys
import os

FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")
transaction_verification_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/transaction_verification'))
sys.path.insert(0, transaction_verification_grpc_path)

import transaction_verification_pb2_grpc as transaction_verification_grpc
import transaction_verification_pb2 as transaction_verification

class TransactionVerificationService(transaction_verification_grpc.TransactionVerificationServiceServicer):
    def checkTransaction(self, request, context):
        request_data = json.loads(request.json)
        # Check if the list of items is not empty
        if not request_data.get('items'):
            response = transaction_verification.TransactionVerificationResponse()
            response.code = "400"
            print("No items in the request")
            return response

        # Check if required user data is filled in
        required_fields = ['user', 'creditCard', 'items', 'billingAddress', 'shippingMethod', 'giftWrapping', 'termsAndConditionsAccepted']
        for field in required_fields:
            if field not in request_data or not request_data[field]:
                response = transaction_verification.TransactionVerificationResponse()
                response.code = "400"
                print(f"Required {field} field not filled in")
                return response

        # If all checks pass, return a successful response
        response = transaction_verification.TransactionVerificationResponse()
        response.code = "200"
        print("Transaction is valid")
        return response




        # Create a TransactionVerificationResponse object
        response = transaction_verification.TransactionVerificationResponse()
        response.code = "200"
        # Set the greeting field of the response object
        # print(request)
        # Return the response object
        return response

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor())
    # Add TransactionDetectionService
    transaction_verification_grpc.add_TransactionVerificationServiceServicer_to_server(TransactionVerificationService(), server)
    # Listen on port 50051
    port = "50052"
    server.add_insecure_port("[::]:" + port)
    # Start the server
    server.start()
    print("Server started. Listening on port 50052.")
    # Keep thread alive
    server.wait_for_termination()

if __name__ == '__main__':
    serve()