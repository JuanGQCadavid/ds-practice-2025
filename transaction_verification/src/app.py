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
        response = transaction_verification.TransactionVerificationResponse()

        if not request_data.get('items'):       
            print("No items in the request")
            response.isValid = False
            response.errMessage = "Missing items field"
            return response


        required_fields = ['user', 'creditCard', 'items', 'billingAddress', 'shippingMethod', 'giftWrapping', 'termsAndConditionsAccepted']
        for field in required_fields:
            if field not in request_data or not request_data[field]:
                print(f"Required {field} field not filled in")
                response.isValid = False
                response.errMessage = f"Missing {field} field"
                return response

        print("Transaction is valid")
        response.isValid = True
        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor())
    transaction_verification_grpc.add_TransactionVerificationServiceServicer_to_server(TransactionVerificationService(), server)
    port = "50052"
    server.add_insecure_port("[::]:" + port)
    server.start()
    print("Server started. Listening on port 50052.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()