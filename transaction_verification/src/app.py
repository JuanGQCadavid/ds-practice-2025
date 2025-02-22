from datetime import datetime
from concurrent import futures
from dotenv import load_dotenv
import google.generativeai as genai
import grpc
import json
import sys
import os

FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")
transaction_verification_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/transaction_verification'))
sys.path.insert(0, transaction_verification_grpc_path)

import transaction_verification_pb2_grpc as transaction_verification_grpc
import transaction_verification_pb2 as transaction_verification

# Load the environment variables
load_dotenv()
genai.configure(api_key=os.environ["GEMINI_API_KEY"])


class TransactionVerificationService(transaction_verification_grpc.TransactionVerificationServiceServicer):
    def checkTransaction(self, request, context):
        print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Checking transaction...")
        request_data = json.loads(request.json)
        response = transaction_verification.TransactionVerificationResponse()
        model = genai.GenerativeModel("gemini-1.5-flash")

        prompt = (
            "Based on the following transaction details, return a risk score from 0 to 99 (just the number), the more, "
            "the most risk:\n"
            f"User: {request_data['user']}\n"
            f"Credit Card: {request_data['creditCard']}\n"
            f"Billing Address: {request_data['billingAddress']}\n"
            f"Shipping Address: {request_data['shippingMethod']}\n"
            f"Items: {request_data['items']}\n"
        )

        response_gemini = model.generate_content(prompt)
        number = int(response_gemini.text.strip())
        print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Transaction validity score: {number}")
        if number > 80:
            response.isValid = False
            response.errMessage = "Not valid transaction"
            return response
        response.isValid = True
        return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor())
    transaction_verification_grpc.add_TransactionVerificationServiceServicer_to_server(TransactionVerificationService(), server)
    port = "50052"
    server.add_insecure_port("[::]:" + port)
    server.start()
    print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Server started. Listening on port 50052.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()