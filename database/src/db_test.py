import grpc
import sys
import os
import time

# Import generated classes
sys.path.insert(0, os.path.abspath(os.path.join(__file__, '../../../utils/pb')))
sys.path.insert(0, os.path.abspath(os.path.join(__file__, '../../../utils/pb/common')))
sys.path.insert(0, os.path.abspath(os.path.join(__file__, '../../../utils/pb/database')))
sys.path.insert(0, os.path.abspath(os.path.join(__file__, '../../../utils/pb/replica')))

import database_pb2 as database
import database_pb2_grpc as database_grpc
import replica_pb2 as replica

REPLICA_PORTS = {
    1: "localhost:50061",
    2: "localhost:50062",
    3: "localhost:50063",
}

global leader_id

def get_stub(replica_id):
    channel = grpc.insecure_channel(REPLICA_PORTS[replica_id])
    return database_grpc.DatabaseServiceStub(channel), channel

def find_leader():
    global leader_id
    for replica_id in REPLICA_PORTS:
        stub, _ = get_stub(replica_id)
        try:
            response = stub.getStatus(replica.StatusRequest(), timeout=2)
            leader_id = response.leaderID
            print(f"[INFO] Leader is {leader_id}")
            return leader_id
        except grpc.RpcError as e:
            print(f"[WARN] Could not contact replica {replica_id}: {e}")
            print("[ERROR] No leader found.")
        return None

def write_stock(book_name, book_stock):
    leader_id = find_leader()
    if leader_id == -1:
        print("[ERROR] Cannot write stock because no valid leader was found.")
        return
    stub, _ = get_stub(leader_id)
    try:
        book_id = str(int(time.time() * 1000))  # Use current time in milliseconds as a unique ID
        response = stub.performStockWrite(
            database.StockRequest(
                bookID=book_id,
                bookName=book_name,
                bookStock=book_stock
            )
        )
    except grpc.RpcError as e:
        print(f"[ERROR] gRPC call to performStockWrite failed: {e}")
        return

    if not response.isValid and response.leaderID:
        new_leader_id = response.leaderID
        print(f"[INFO] Forwarding write request to new leader: {new_leader_id}")
        stub_new, _ = get_stub(new_leader_id)
        try:
            book_id = str(int(time.time() * 1000))  # Use current time in milliseconds as a unique ID
            response = stub_new.performStockWrite(
                database.StockRequest(
                    bookID=book_id,
                    bookName=book_name,
                    bookStock=book_stock
                )
            )
        except grpc.RpcError as e:
            print(f"[ERROR] gRPC call to new leader {new_leader_id} failed: {e}")
            return

    if response.isValid:
        print(f"[SUCCESS] Stock updated {book_id}: {book_name} = {book_stock}")
    else:
        print(f"[ERROR] {response.errMessage}")

def read_stock(book_id):
    leader_id = find_leader()

    stub, _ = get_stub(leader_id)
    try:
        response = stub.performStockRead(database.StockRequest(bookID=book_id))
        if response.isValid:
            print(f"[INFO] Replica {leader_id} reports {book_id} = {response.bookName, response.bookStock}")
        else:
            print(f"[INFO] Replica {leader_id} - {response.errMessage}")
    except grpc.RpcError as e:
        print(f"[WARN] Could not contact replica {leader_id}: {e}")

if __name__ == '__main__':
    print("Starting client...")
    while True:
        action = input("\nChoose action: [w]rite / [r]ead / [q]uit: ").strip().lower()
        if action == 'q':
            print("Exiting client.")
            break
        elif action == 'w':
            book = input("Enter book name: ")
            try:
                qty = int(input("Enter stock quantity: "))
            except ValueError:
                print("[ERROR] Quantity must be an integer.")
                continue
            write_stock(book, qty)
        elif action == 'r':
            book_id = input("Enter book ID to read: ")
            read_stock(book_id)
        else:
            print("Invalid option.")

