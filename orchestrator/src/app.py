import sys
import os
from concurrent import futures
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

def greet(name='you'):
    # Establish a connection with the fraud-detection gRPC service.
    with grpc.insecure_channel('fraud_detection:50051') as channel:
        # Create a stub object.
        stub = fraud_detection_grpc.HelloServiceStub(channel)
        # Call the service through the stub object.
        response = stub.SayHello(fraud_detection.HelloRequest(name=name))
    return response.greeting

def check_transaction(request):
    request_json = json.dumps(request)
    # Establish connection with the transaction_verification gRPC servive.
    with grpc.insecure_channel('transaction_verification:50052') as channel:
        stub = transaction_verification_grpc.TransactionVerificationServiceStub(channel)
        # Call the service through the stub object
        response = stub.checkTransaction(transaction_verification.TransactionVerificationRequest(json=request_json))
    return response.code

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
    response = greet(name='orchestrator')
    # Return the response.
    return response

@app.route('/checkout', methods=['POST'])
def checkout():
    """
    Responds with a JSON object containing the order ID, status, and suggested books.
    """
    # Get request object data to json
    request_data = json.loads(request.data)
    # Print request object data
    print("Request Data:", request_data.get('items')) # on terminal

    response = None

    # TODO: create microservices
    with futures.ThreadPoolExecutor() as executor:
        # Check fraud
        future_fraud_detection = executor.submit(check_fraud, request_data)
        # Check transaction
        future_transaction_verification = executor.submit(check_transaction, request_data)
        # Get suggestions
        future_suggestions = executor.submit(get_suggestions, request_data)

        fraud_detection_code = future_fraud_detection.result()
        transaction_verification_code = future_transaction_verification.result()
        suggestions = future_suggestions.result()
    transaction_verification_code = check_transaction(request_data)
    if fraud_detection_code != "200":
        # set error message
        error_response = {
            'error': {
                'code': '400',
                'message': 'Fraud detected'
            }
        }
        return error_response, 400
    elif transaction_verification_code != "200":
        # set error message
        response = {
            'error': {
                'code': '400',
                'message': 'Transaction is invalid'
            }
        }
        return response, 400
    response = {
        'orderId': '12345',
        'status': 'Order Approved',
        'suggestedBooks': suggestions
    }
    return response


if __name__ == '__main__':
    # Run the app in debug mode to enable hot reloading.
    # This is useful for development.
    # The default port is 5000.
    app.run(host='0.0.0.0')
