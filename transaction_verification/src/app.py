import time
from datetime import datetime
from concurrent import futures
from dotenv import load_dotenv
import google.generativeai as genai
import grpc
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
    def __init__(self, service_index=0, total_services=3):
        self.service_index = service_index
        self.max_services = total_services
        self.orders = {}
        self.ai_model = genai.GenerativeModel("gemini-1.5-flash")
        self.vector_clock_access = True  # variable to check vc access on concurrent execution

    def initOrder(self,request, context):
        order_id = request.orderId
        order_data = request.order

        print(order_id)
        print(order_data)
        response = transaction_verification.TransactionVerificationResponse()

        if order_id not in self.orders:
            self.orders[order_id] = {
                "fields": order_data,
                "vc":  [0] * self.max_services
            }
            response.isValid = True
        else:
            response.isValid = False
            response.errMessage = "OrderID already exists"
        return response

    def merge_and_increment(self, local_vc, received_vc):
        while not self.vector_clock_access: # critical section usage
            continue # dummy wait
        self.vector_clock_access = False
        for i in range(self.max_services):
            local_vc[i] = max(local_vc[i], received_vc[i])
        local_vc[self.service_index] += 1
        self.vector_clock_access = True

    def checkOrder(self, request, context):
        order_id = request.orderId 
        incoming_vc = request.clock

        # print(self.orders)

        entry = self.orders[order_id]
        print(entry)

        self.merge_and_increment(entry["vc"], incoming_vc)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_order {entry['vc']}")

        hola = transaction_verification.TransactionVerificationResponseClock()
        standard_response = transaction_verification.TransactionVerificationResponse()
        
        print("--------------------------------------")
        print(type(entry["fields"]) )
        print(entry.get("user"))

        if (not entry["fields"].user or not entry["fields"].creditCard or not entry["fields"].items or not entry["fields"].shippingMethod or
                not entry["fields"].billingAddress):
            standard_response.isValid = False
            standard_response.errMessage = "No data"
            # response.clock = entry["vc"]
            # hola.response = standard_response
            return transaction_verification.TransactionVerificationResponseClock(
                response=standard_response,
            )
    

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_order completed")
        standard_response.isValid = True
        response.clock = entry["vc"]
        return response

    def check_user(self, order_id, incoming_vc):
        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], incoming_vc)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_user {entry['vc']}")

        order_data = entry["order"]
        prompt = (
            "Based on the user data (name, contact, address), return a risk score from 0 to 99 (just the number), "
            "the highest, the most untruthful:\n"
            f"User: {order_data['user']}\n"
            f"Billing address: {order_data['billing_address']}\n"
            f"Shipping address: {order_data['shipping_address']}\n"
        )

        response_ai_model = self.ai_model.generate_content(prompt)
        user_score = int(response_ai_model.text.strip())
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} check_user score {user_score}")

        response = transaction_verification.TransactionVerificationResponseClock()
        standard_response = transaction_verification.TransactionVerificationResponse()
        if user_score > 80:
            standard_response.isValid = False
            standard_response.message = "User score is higher than 80"
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_user completed")
        standard_response.isValid = True
        response.clock = entry["vc"]
        return response

    def check_format_credit_card(self, order_id, incoming_vc):
        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], incoming_vc)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_format_credit_card {entry['vc']}")

        credit_card = entry["credit_card"]
        prompt = (
            "Based on the following credit card, check the format. If it is right, return a 1, else return 0\n"
            f"Credit card: {credit_card}\n"
        )

        response_ai_model = self.ai_model.generate_content(prompt)
        format_validity = int(response_ai_model.text.strip())

        response = transaction_verification.TransactionVerificationResponseClock()
        standard_response = transaction_verification.TransactionVerificationResponse()

        if not format_validity:
            standard_response.isValid = False
            standard_response.message = "Format credit card is invalid"

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} credit card check successfully completed")
        standard_response.isValid = True
        response.clock = entry["vc"]
        return response

    # def checkTransaction(self, request, context):
    #     print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Checking transaction...")
    #     request_data = json.loads(request.json)
    #     response = transaction_verification.TransactionVerificationResponse()
    #     model = genai.GenerativeModel("gemini-1.5-flash")
    #
    #     prompt = (
    #         "Based on the following transaction details, return a risk score from 0 to 99 (just the number), the more, "
    #         "the most risk:\n"
    #         f"User: {request_data['user']}\n"
    #         f"Credit Card: {request_data['creditCard']}\n"
    #         f"Billing Address: {request_data['billingAddress']}\n"
    #         f"Shipping Address: {request_data['shippingMethod']}\n"
    #         f"Items: {request_data['items']}\n"
    #     )
    #
    #     response_gemini = model.generate_content(prompt)
    #     number = int(response_gemini.text.strip())
    #     print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Transaction validity score: {number}")
    #     if number > 80:
    #         response.isValid = False
    #         response.errMessage = "Not valid transaction"
    #         return response
    #     response.isValid = True
    #     return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor())
    transaction_verification_grpc.add_TransactionVerificationServiceServicer_to_server(TransactionVerificationService(), server)
    port = "50052"
    server.add_insecure_port("[::]:" + port)
    server.start()
    print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Server started. Listening on port 50052.")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()