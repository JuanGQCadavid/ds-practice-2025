from datetime import datetime
from concurrent import futures
from dotenv import load_dotenv
from google.protobuf import json_format
import google.generativeai as genai
import sys
import os

FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")

all_pb = os.path.abspath(os.path.join(FILE, '../../../utils/pb'))
sys.path.insert(0, all_pb)

common_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/common'))
sys.path.insert(0, common_grpc_path)

transaction_verification_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/transaction_verification'))
sys.path.insert(0, transaction_verification_grpc_path)

import transaction_verification_pb2_grpc as transaction_verification_grpc
import transaction_verification_pb2 as transaction_verification
import common_pb2 as common_pb
import grpc
from concurrent import futures

# Load the environment variables
load_dotenv()
genai.configure(api_key=os.environ["GEMINI_API_KEY"])

class TransactionVerificationService(transaction_verification_grpc.TransactionVerificationServiceServicer):
    def __init__(self, service_index=0, total_services=3):
        self.service_index = service_index
        self.max_services = total_services
        self.orders = {}
        self.ai_model = genai.GenerativeModel("gemini-1.5-flash-002")
        self.vector_clock_access = True  # variable to check vc access on concurrent execution

    def initOrder(self, request, context):
        order_id = request.orderId
        order = request.order

        order_json = json_format.MessageToDict(order)
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Received order ID {order_id}")

        response = common_pb.InitResponse()
        
        if order_id in self.orders:
            response.isValid = False
            response.errMessage = "Order ID already exists"
            return response

        self.orders[order_id] = {
            "user": order.user,
            "credit_card": order.creditCard,
            "user_comment": order.userComment,
            "items": order.items,
            "discount_code": order.discountCode,
            "shipping_method": order.shippingMethod,
            "gift_message": order.giftMessage,
            "billing_address": order.billingAddress,
            "gift_wrapping": order.giftWrapping,
            "terms_accepted": order.termsAccepted,
            "notification_preferences": order.notificationPreferences,
            "device": order.device,
            "browser": order.browser,
            "app_version": order.appVersion,
            "screen_resolution": order.screenResolution,
            "referrer": order.referrer,
            "device_language": order.deviceLanguage,
            "vc": [0] * self.max_services,
        }
        response.isValid = True
        return response

    def merge_and_increment(self, local_vc, received_vc):
        while not self.vector_clock_access: # critical section usage
            continue # dummy wait
        self.vector_clock_access = False
        for i in range(self.max_services):
            local_vc[i] = max(local_vc[i], received_vc[i])
        local_vc[self.service_index] += 1
        self.vector_clock_access = True

    def exist_order(self, order_id):
        if order_id not in self.orders:
            return False
        return True

    def error_response(self, message):
        return common_pb.NextResponse(
            errMessage = message,
            isValid = False,
        )

    def checkOrder(self, request, context):
        order_id = request.orderId 

        if not self.exist_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID{order_id} does not exist")
            return self.error_response("Order does not exist")

        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], request.incomingVectorClock)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkOrder {entry['vc']}")

        if (not entry["user"] or not entry["credit_card"] or not entry["items"] or not entry["shipping_method"] or
                not entry["billing_address"]):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Invalid order: some fields are empty")
            return self.error_response("Invalid order: some fields are empty")

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkOrder completed")
        return common_pb.NextResponse(
            vectorClock = entry["vc"],
            isValid = True,
        )

    def checkUser(self, request, context):
        order_id = request.orderId

        if not self.exist_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} does not exist")
            return self.error_response("Order does not exist")

        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], request.incomingVectorClock)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkUser {entry['vc']}")

        prompt = (
            "Based on the user data (name, contact, address), return a risk score from 0 to 99 (just the number), "
            "the highest, the most untruthful:\n"
            f"User: {entry['user']}"
            f"Billing address: {entry['billing_address']}"
        )

        response_ai_model = self.ai_model.generate_content(prompt)
        user_score = int(response_ai_model.text.strip())

        if user_score > 80:
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} User data is suspicious")
            return self.error_response("User data is suspicious")
        
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkUser completed")

        return common_pb.NextResponse(
            vectorClock = entry["vc"],
            isValid = True,
        )

    def checkFormatCreditCard(self, request, context):
        order_id = request.orderId

        if not self.exist_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} does not exist")
            return self.error_response("Order does not exist")

        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], request.incomingVectorClock)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkFormatCreditCard {entry['vc']}")
        
        if len(entry["credit_card"].number) != 16:  # strict to VISA / MasterCard
            print(entry["credit_card"].number)
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Credit card number format is not valid")
            return self.error_response("Credit card number format is not valid")

        if len(entry["credit_card"].cvv) != 3:  # strict to VISA / MasterCard format
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} CVV number format is not valid")
            return self.error_response("CVV number format is not valid")

        if (not entry["credit_card"].expirationDate or not entry["credit_card"].expirationDate[0:2].isdigit() or
                not entry["credit_card"].expirationDate[3:5].isdigit()):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Expiration date format is not valid")
            return self.error_response("Expiration date format is not valid")

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkFormatCreditCard completed")
        return common_pb.NextResponse(
            vectorClock = entry["vc"],
            isValid = True,
        )
        
    def cleanOrder(self, request, context):
        order_id = request.orderId
        
        if not self.exist_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} does not exist")
            return self.error_response("Order does not exist")
        
        self.orders.pop(order_id)
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} removed from pending")
        return common_pb.NextResponse(
            isValid = True,
        )
        
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