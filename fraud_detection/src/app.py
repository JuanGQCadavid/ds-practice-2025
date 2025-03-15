import sys
import os
from datetime import datetime
from google.protobuf import json_format

# This set of lines are needed to import the gRPC stubs.
# The path of the stubs is relative to the current file, or absolute inside the container.
# Change these lines only if strictly needed.
FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")
fraud_detection_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/fraud_detection'))
sys.path.insert(0, fraud_detection_grpc_path)
import fraud_detection_pb2 as fraud_detection
import fraud_detection_pb2_grpc as fraud_detection_grpc
import grpc
from concurrent import futures

class FraudDetectionService(fraud_detection_grpc.FraudDetectionServiceServicer):
    def __init__(self, service_index=1, total_services=3):
        self.service_index = service_index
        self.max_services = total_services
        self.orders = {}
        self.vector_clock_access = True

    def initOrder(self, request, context):
        order_id = request.orderId
        order = request.order

        order_json = json_format.MessageToDict(order)
        print(f"Received order ID {order_id}")
        print(f"{order_json}")

        response = fraud_detection.FraudDetectionResponse()

        if order_id not in self.orders:
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
        else:
            response.isValid = False
            response.errMessage = "OrderID already exists"
        return response

    def merge_and_increment(self, local_vc, received_vc):
        while not self.vector_clock_access:
            continue
        self.vector_clock_access = False
        for i in range(self.max_services):
            local_vc[i] = max(local_vc[i], received_vc[i])
        local_vc[self.service_index] += 1
        self.vector_clock_access = True

    def exists_order(self, order_id):
        if order_id not in self.orders:
            return False
        return True

    def error_response(self):
        return fraud_detection.FraudDetectionResponseClock(
            response = fraud_detection.FraudDetectionResponse(
                code = "200",
            )
        )

    def checkUser(self, request, context):
        order_id = request.orderId

        if not self.exists_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} does not exist")
            return self.error_response()

        incoming_vc = request.clock
        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], incoming_vc)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkUser {entry['vc']}")

        user_data = entry["user"]
        name = entry["user"].name
        contact = entry["user"].contact

        suspicious_names = ["John Doe", "Jane Doe", "Test User"]
        suspicious_contact = ["john.doe@example.com", "jane.doe@example.com", "test.user@example.com"]

        if name in suspicious_names or contact in suspicious_contact:
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkUser: User is suspicious")
            return fraud_detection.FraudDetectionResponseClock(
                response = fraud_detection.FraudDetectionResponse(
                    code = "400",
                ),
            )

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} checkUser completed")
        return fraud_detection.FraudDetectionResponseClock(
            response = fraud_detection.FraudDetectionResponse(
                code = "200",
            ),
            clock = entry["vc"],
        )

    def checkCreditCard(self, request, context):
        order_id = request.orderId

        if not self.exists_order(order_id):
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} does not exist")
            return self.error_response()

        incoming_vc = request.clock
        entry = self.orders[order_id]
        self.merge_and_increment(entry["vc"], incoming_vc)

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} checkCreditCard {entry['vc']}")

        card_number = entry["credit_card"].number
        cvv = entry["credit_card"].cvv
        expiration_date = entry["credit_card"].expirationDate

        suspicious_card_numbers = ["4111111111111111", "5500000000000004", "340000000000009"]
        suspicious_cvv = ["000", "111", "123"]

        if card_number in suspicious_card_numbers or cvv in suspicious_cvv:
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_credit_card: Credit card is suspicious")
            return self.error_response()

        expiration_date = datetime.strptime(expiration_date, "%m/%y")
        if expiration_date < datetime.now():
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_credit_card: Credit card is expired")
            return self.error_response()

        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order {order_id} check_credit_card completed")
        return fraud_detection.FraudDetectionResponseClock(
            response = fraud_detection.FraudDetectionResponse(
                code = "200",
            ),
            clock = entry["vc"],
        )

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor())
    # Add FraudDetectionService to the server
    fraud_detection_grpc.add_FraudDetectionServiceServicer_to_server(FraudDetectionService(), server)
    # Listen on port 50051
    port = "50051"
    server.add_insecure_port("[::]:" + port)
    # Start the server
    server.start()
    print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Server started on port {port}")
    # Keep thread alive
    server.wait_for_termination()

if __name__ == '__main__':
    serve()