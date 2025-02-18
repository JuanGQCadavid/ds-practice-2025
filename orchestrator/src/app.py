import sys
import os
from concurrent import futures
from datetime import datetime
from http.client import responses

# This set of lines are needed to import the gRPC stubs.
# The path of the stubs is relative to the current file, or absolute inside the container.
# Change these lines only if strictly needed.
FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")
fraud_detection_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/fraud_detection'))
sys.path.insert(0, fraud_detection_grpc_path)
import fraud_detection_pb2 as fraud_detection
import fraud_detection_pb2_grpc as fraud_detection_grpc

transaction_verification_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/transaction_verification'))
sys.path.insert(0, transaction_verification_grpc_path)
import transaction_verification_pb2 as transaction_verification
import transaction_verification_pb2_grpc as transaction_verification_grpc

import grpc

def check_fraud(request):
    card_data = request.get('creditCard')
    if not card_data: # Check if creditCard key is present in the request
        return "400"

    # Encapsulate in gRPC request object
    credit_card = fraud_detection.CreditCard(
        number=card_data.get('number'),
        cvv=card_data.get('cvv'),
        expirationDate=card_data.get('expirationDate')
    )
    fraud_request = fraud_detection.FraudDetectionRequest(creditCard=credit_card)

    # Make gRPC call
    with grpc.insecure_channel('fraud_detection:50051') as channel:
        stub = fraud_detection_grpc.FraudDetectionServiceStub(channel)
        response = stub.checkFraud(fraud_request)  # Pass the request directly
    return response.code

def check_transaction(request):
    request_json = json.dumps(request)
    # Establish connection with the transaction_verification gRPC service.
    with grpc.insecure_channel('transaction_verification:50052') as channel:
        stub = transaction_verification_grpc.TransactionVerificationServiceStub(channel)
        # Call the service through the stub object
        response = stub.checkTransaction(transaction_verification.TransactionVerificationRequest(json=request_json))
    return response

# Import Flask.
# Flask is a web framework for Python.
# It allows you to build a web application quickly.
# For more information, see https://flask.palletsprojects.com/en/latest/
from flask import Flask, request
from flask_cors import CORS
import json

# Create a simple Flask app.
app = Flask(__name__)
# Enable CORS for the app.
CORS(app, resources={r'/*': {'origins': '*'}})

# Define a GET endpoint.
@app.route('/', methods=['GET'])
def index():
    """
    Responds with 'Hello, [name]' when a GET request is made to '/' endpoint.
    """
    # Test the fraud-detection gRPC service.
    response = "I am the orchestrator"
    # Return the response.
    return response

@app.route('/checkout', methods=['POST'])
def checkout():
    """
    Responds with a JSON object containing the order ID, status, and suggested books.
    """
    # Get request object data to json
    request_data = json.loads(request.data)

    fraud_detection_code = None  # initialize fraud_detection_code
    transaction_verification_code = None  # initialize transaction_verification_code
    suggestions = None  # initialize suggestions

    with futures.ThreadPoolExecutor() as executor:
        # Check fraud
        future_fraud_detection = executor.submit(check_fraud, request_data)
        # Check transaction
        future_transaction_verification = executor.submit(check_transaction, request_data)
        # Get suggestions
        #future_suggestions = executor.submit(get_suggestions, request_data)

        fraud_detection_code = future_fraud_detection.result()
        transaction_verification_response = future_transaction_verification.result()
        #suggestions = future_suggestions.result()
    if fraud_detection_code != "200":
        # set error message
        error_response = {
            'error': {
                'code': '400',
                'message': 'Fraud detected'
            }
        }
        return error_response, 400
    elif transaction_verification_response["isValid"] != "true":
        # set error message
        response = transaction_verification_response["message"]
        return response, 400
    response = {
        'orderId': '12345',
        'status': 'Order Approved',
        'suggestedBooks': [
            {'bookId': '123', 'title': 'The Best Book', 'author': 'Author 1'},
            {'bookId': '456', 'title': 'The Second Best Book', 'author': 'Author 2'}
        ]
    }
    return response


if __name__ == '__main__':
    # Run the app in debug mode to enable hot reloading.
    # This is useful for development.
    # The default port is 5000.
    # add logs with timestamp
    print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Orchestrator service started")
    app.run(host='0.0.0.0')
