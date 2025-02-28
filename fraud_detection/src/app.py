import sys
import os
from datetime import datetime

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

# Create a class to define the server functions, derived from
# fraud_detection_pb2_grpc.HelloServiceServicer
class FraudDetectionService(fraud_detection_grpc.FraudDetectionServiceServicer):
    def validate_card_details(self, card_number, cvv, expiration_date):
        # Card Number Validation
        def validate_card_number(card_number):
            # Luhn Algorithm for card number validation
            def luhn_checksum(card_num):
                digits = [int(d) for d in str(card_num)]
                checksum = 0
                is_even = False
                for digit in reversed(digits):
                    if is_even:
                        digit *= 2
                        if digit > 9:
                            digit -= 9
                    checksum += digit
                    is_even = not is_even
                return checksum % 10 == 0
            return (len(str(card_number)) in [13, 14, 15, 16]) and luhn_checksum(card_number)
        # CVV Validation
        def validate_cvv(cvv):
            return len(str(cvv)) in [3, 4]
        # Expiration Date Validation
        def validate_expiration_date(exp_date):
            import datetime
            try:
                # Assuming format is MM/YY
                month, year = map(int, exp_date.split('/'))
                current_date = datetime.datetime.now()
                exp_date = datetime.datetime(2000 + year, month, 1)
                return exp_date > current_date
            except:
                return False
        # Combine validations
        results = {
            'card_number_valid': validate_card_number(card_number),
            'cvv_valid': validate_cvv(cvv),
            'expiration_valid': validate_expiration_date(expiration_date)
        }
        return results

    def assess_fraud_risk(self, card_number, cvv, validation_results):
        risk_factors = {
            'high_risk_bins': ['4111', '5500', '3400'],  # Example BIN numbers
            'suspicious_cvv_patterns': ['000', '111', '123']
        }
        risk_score = 0

        # Check BIN (first 4 digits)
        if str(card_number)[:4] in risk_factors['high_risk_bins']:
            risk_score += 30

        # Check CVV pattern
        if str(cvv) in risk_factors['suspicious_cvv_patterns']:
            risk_score += 20

        # Validation check
        if not all(validation_results.values()): # If any of the values is False
            risk_score += 50
        print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Risk Score:", risk_score)
        # Categorize risk
        if risk_score > 50:
            return "High Risk"
        elif risk_score > 20:
            return "Medium Risk"
        else:
            return "Low Risk"

    # Create an RPC function to say checkFraud
    def checkFraud(self, request, context):
        print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] Checking fraud...")
        # Access credit card fields directly from the request
        credit_card = request.creditCard
        validation_results = self.validate_card_details(
            credit_card.number,
            credit_card.cvv,
            credit_card.expirationDate
        )
        risk_assessment = self.assess_fraud_risk(
            credit_card.number,
            credit_card.cvv,
            validation_results
        )
        # Create and return response
        response = fraud_detection.FraudDetectionResponse()
        if risk_assessment == "High Risk":
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} High risk detected, fraud suspected, invalidating purchase")
            response.code = "400"
        else:
            response.code = "200"
        return response

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor())
    # Add HelloService
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