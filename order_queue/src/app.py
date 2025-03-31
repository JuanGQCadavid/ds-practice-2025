import sys
import os
from datetime import datetime
from google.protobuf import json_format
import heapq
import threading

# This set of lines are needed to import the gRPC stubs.
# The path of the stubs is relative to the current file, or absolute inside the container.
# Change these lines only if strictly needed.
FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")

all_pb = os.path.abspath(os.path.join(FILE, '../../../utils/pb'))
sys.path.insert(0, all_pb)

common_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/common'))
sys.path.insert(0, common_grpc_path)

order_queue_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/order_queue'))
sys.path.insert(0, order_queue_grpc_path)

import order_queue_pb2 as order_queue
import order_queue_pb2_grpc as order_queue_grpc
import common_pb2 as common_pb
import grpc
from concurrent import futures

class OrderQueueService(order_queue_grpc.OrderQueueServiceServicer):
    def __init__(self):
        self._lock = threading.Lock()
        self._queue = []

    def enqueue(self, request, context):
        # lock the queue
        response = common_pb.NextResponse()
        response.isValid = True
        with self._lock:
            priority = 3 # default priority
            order = request.order
            if order.clientCard == "None":
                priority = 2
            elif order.clientCard == "Basic":
                priority = 1
            elif order.clientCard == "Premium":
                priority = 0
            else:
                response.isValid = False
                response.errMessage = "clientCard field not among the allowed values"
            heapq.heappush(self._queue, (priority, request.orderId, order))
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {request.orderId} enqueued")
            return response

    def dequeue(self, request, context):
        with self._lock:
            if len(self._queue) == 0:
                print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Queue is empty, request not served")
                return order_queue.DequeueResponse(
                    isValid=False,
                    errMessage="Queue is empty"
                )
            _, order_id, order = self._queue[0]
            heapq.heappop(self._queue)
            print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Order ID {order_id} dequeued")
            return order_queue.DequeueResponse(
                orderId = order_id,
                order = order,
                isValid = True,
            )

    def clean(self, request, context):
        with self._lock:
            self._queue = []
        return common_pb.NextResponse(isValid=True)

def serve():
    # Create a gRPC server
    server = grpc.server(futures.ThreadPoolExecutor())
    # Add OrderQueueService to the server
    order_queue_grpc.add_OrderQueueServiceServicer_to_server(OrderQueueService(), server)
    # Listen on port 50054
    port = "50054"
    server.add_insecure_port("[::]:" + port)
    # Start the server
    server.start()
    print(f"{datetime.now().strftime('%Y/%m/%d %H:%M:%S')} Server started on port {port}")
    # Keep thread alive
    server.wait_for_termination()

if __name__ == '__main__':
    serve()