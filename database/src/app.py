import sys
import os
import threading
import time
from datetime import datetime

# This set of lines are needed to import the gRPC stubs.
# The path of the stubs is relative to the current file, or absolute inside the container.
# Change these lines only if strictly needed.
FILE = __file__ if '__file__' in globals() else os.getenv("PYTHONFILE", "")

all_pb = os.path.abspath(os.path.join(FILE, '../../../utils/pb'))
sys.path.insert(0, all_pb)

common_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/common'))
sys.path.insert(0, common_grpc_path)

database_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/database'))
sys.path.insert(0, database_grpc_path)

replica_grpc_path = os.path.abspath(os.path.join(FILE, '../../../utils/pb/replica'))
sys.path.insert(0, replica_grpc_path)

import database_pb2 as database
import database_pb2_grpc as database_grpc
import common_pb2 as common_pb
import replica_pb2 as replica
import replica_pb2_grpc as replica_grpc
import grpc
from concurrent import futures

class DatabaseService(database_grpc.DatabaseServiceServicer, replica_grpc.ReplicaServicer):
    def __init__(self):
        self.rank = int(os.environ.get("RANK", -1))
        self.stock = {}
        self.last_access = datetime.now().strftime('%Y/%m/%d %H:%M:%S')
        self.replicas_id = {self.rank}  # Incluir nuestro propio ID desde el inicio
        self.replica_endpoints = {
            1: f"db_replica1:50061",
            2: f"db_replica2:50062",
            3: f"db_replica3:50063",
        }
        self.num_replicas = 3  # static number of replicas
        self.leader_rank = -1  # -1 means no leader
        self.lock = threading.RLock()  # Void race conditions

        # Variables para heartbeat
        self.last_heartbeat = time.time()
        self.heartbeat_interval = 1
        self.heartbeat_timeout = 5
        self.heartbeat_thread = None
        self.timeout_thread = None

        self._channels = {}
        self._stubs = {}

        # Start leader election on boot
        threading.Thread(target=self._delayed_election_start, daemon=True).start()

        self.timeout_thread = threading.Thread(target=self.check_leader_timeout, daemon=True)
        self.timeout_thread.start()

    def _delayed_election_start(self):
        time.sleep(2)
        self.log("Broadcasting my ID to all replicas")
        self.broadcast("shareID", replica.IDRequest(ID=self.rank), exclude_self=False)
        self.log("Automatically starting leader election")
        self.start_election_process()

    def log(self, message):
        leader_status = "LEADER" if self.leader_rank == self.rank else f"FOLLOWER (Leader: {self.leader_rank})"
        timestamp = datetime.now().strftime('%Y/%m/%d %H:%M:%S')
        print(f"{timestamp} [{leader_status}]: {message}")

    def get_stub(self, rank):
        endpoint = self.replica_endpoints.get(rank)
        if not endpoint:
            return None

        if rank in self._channels:
            try:
                grpc.channel_ready_future(self._channels[rank]).result(timeout=1)
            except Exception as e:
                self._channels.pop(rank)
                self._stubs.pop(rank, None)

        if rank not in self._channels:
            self._channels[rank] = grpc.insecure_channel(endpoint)
            self._stubs[rank] = replica_grpc.ReplicaStub(self._channels[rank])

        return self._stubs[rank]

    def broadcast(self, method_name, request, exclude_self=True):
        responses = {}

        for rep_rank in self.replica_endpoints:
            if exclude_self and rep_rank == self.rank:
                continue
            try:
                stub = self.get_stub(rep_rank)
                if stub:
                    method = getattr(stub, method_name)
                    response = method(request, timeout=3)
                    responses[rep_rank] = response
                    self.log(f"Message {method_name} sent to replica {rep_rank}")
            except Exception as e:
                self.log(f"Error sending {method_name} to replica {rep_rank}")
                responses[rep_rank] = None
        return responses

    def start_election_process(self):
        self.log("Starting leader election")

        responses = self.broadcast("shareID", replica.IDRequest(ID=self.rank), exclude_self=False)
        with self.lock:
            for rep, response in responses.items():
                if response and response.knownIDs:
                    self.replicas_id.update(response.knownIDs)

        time.sleep(3)
        self.elect_leader()

    def elect_leader(self):
        with self.lock:
            if not self.replicas_id:
                self.log("No replicas")
                return

            highest_id = max(self.replicas_id)
            self.log(f"Known IDs: {self.replicas_id}, highest ID: {highest_id}")

            if highest_id == self.rank:
                self.log(f"I'm the highest ({self.rank}), becoming leader")
                self.become_leader()
            else:
                self.log(f"Replica {highest_id} has the highest ID, becoming follower")
                self.become_follower(highest_id)

    def become_leader(self):
        with self.lock:
            old_leader = self.leader_rank
            self.leader_rank = self.rank

            if old_leader != self.rank:
                self.log("YES, I'M THE LEADER!")

                self.broadcast(
                    "notifyNewLeader",
                    replica.LeaderNotification(leaderID=self.rank)
                )

                # Start heartbeat thread if not already running
                if self.heartbeat_thread is None or not self.heartbeat_thread.is_alive():
                    self.heartbeat_thread = threading.Thread(target=self.send_heartbeats, daemon=True)
                    self.heartbeat_thread.start()

    def become_follower(self, leader_id):
        with self.lock:
            old_leader = self.leader_rank
            self.leader_rank = leader_id

            if old_leader != leader_id:
                self.log(f"I'm a FOLLOWER of leader {leader_id}")
                self.last_heartbeat = time.time()  # Reset heartbeat time

    def shareID(self, request, context):
        replica_id = request.ID

        with self.lock:
            if replica_id not in self.replicas_id:
                self.replicas_id.add(replica_id)
                self.log(f"Replica {replica_id} exists. Known IDs: {sorted(self.replicas_id)}")

                if self.leader_rank == self.rank and replica_id > self.rank:  # I am the leader but gotten a higher ID
                    self.log(f"Replica ({replica_id}) took my crown. Well folks, it's time for me to disappear like your motivation after lunch. Remember me not for what I did, but for all the snacks I stole from the break room. Goodbye, and may your coffee be strong and your Wi-Fi never drop!")
                    self.transfer_leadership(replica_id)
            else:
                self.log(f"Replica {replica_id} already exists in the list.")

        response = replica.IDResponse(isValid=True)
        response.knownIDs.extend(sorted(list(self.replicas_id)))
        return response

    def notifyNewLeader(self, request, context):
        leader_id = request.leaderID
        self.log(f"Got replica {leader_id} leader notification")

        self.become_follower(leader_id)

        return replica.IDResponse(isValid=True)

    def transfer_leadership(self, new_leader_id):
        with self.lock:
            if self.leader_rank != self.rank:
                return

            self.log(f"You're new leader is replica {new_leader_id}")

            try:
                stub = self.get_stub(new_leader_id)
                if stub:
                    stub.becomeLeader(replica.LeaderTransfer(
                        oldLeaderID=self.rank,
                        newLeaderID=new_leader_id
                    ))
            except Exception as e:
                self.log(f"Lost contact on leadership notification, how sad")
                return  # If communication fails, we don't need to do anything

            self.broadcast(
                "notifyNewLeader",
                replica.LeaderNotification(leaderID=new_leader_id),
                exclude_self=False
            )

            self.become_follower(new_leader_id)

    def becomeLeader(self, request, context):
        old_leader_id = request.oldLeaderID
        new_leader_id = request.newLeaderID

        self.log(f"Yes, my time has come, I'm the leader of {old_leader_id}")

        if new_leader_id == self.rank:
            self.become_leader()

        return replica.IDResponse(isValid=True)

    def heartbeat(self, request, context):
        leader_id = request.leaderID

        with self.lock:
            if self.leader_rank == leader_id:
                self.last_heartbeat = time.time()
                self.log(f"Replica {leader_id} still alive")

        return replica.HeartbeatResponse(
            replicaID=self.rank,
            isAlive=True
        )

    def send_heartbeats(self):
        self.log("Starting heartbeat communication")
        while True:
            with self.lock:
                if self.leader_rank != self.rank:
                    self.log("It was a pleasure to be your leader")
                    break

                self.log("Still alive, sending heartbeats")
                self.broadcast(
                    "heartbeat",
                    replica.HeartbeatRequest(leaderID=self.rank)
                )

            time.sleep(self.heartbeat_interval)

    def check_leader_timeout(self):
        while True:
            with self.lock:
                if self.leader_rank != self.rank and self.leader_rank != -1:
                    time_since_last = time.time() - self.last_heartbeat
                    if time_since_last > self.heartbeat_timeout:
                        self.log(
                            f"Timeout, leader passed away {self.leader_rank}! Last heartbeat {time_since_last:.2f}s")
                        self.replicas_id.discard(self.leader_rank)
                        self.leader_rank = -1  # Reset leader
                        time.sleep(1)
                        threading.Thread(target=self.start_election_process, daemon=True).start()
            time.sleep(1)

    def performStockRead(self, request, context):
        self.log(f"Read operation received for book: {request.bookID}")

        with self.lock:
            self.last_access = datetime.now().strftime('%Y/%m/%d %H:%M:%S')
            print(self.stock.keys())
            if request.bookID in self.stock.keys():
                book_name, book_stock = self.stock[request.bookID]
                return database.StockResponse(
                    bookName=book_name,
                    bookStock=book_stock,
                    isValid=True
                )
            else:
                return database.StockResponse(
                    bookStock=0,
                    isValid=False,
                    errMessage=f"Book {request.bookID} not found in stock"
                )

    def performStockWrite(self, request, context):
        self.log(f"Write operation received for bookId {request.bookID}: {request.bookName}, quantity: {request.bookStock}")

        with self.lock:
            self.last_access = datetime.now().strftime('%Y/%m/%d %H:%M:%S')

            if self.leader_rank != self.rank:
                self.log(f"Write operation rejected: I'm not the leader. Leader is: {self.leader_rank}")
                return database.StockResponse(
                    bookStock=0,
                    isValid=False,
                    errMessage=f"Not the leader. Forward request to replica {self.leader_rank}",
                    leaderID=self.leader_rank
                )

            self.stock[request.bookID] = [request.bookName, request.bookStock]

            self.propagate_write(request.bookID, request.bookName, request.bookStock)

            return database.StockResponse(
                bookStock=request.bookStock,
                isValid=True
            )

    def propagate_write(self, book_id, book_name, book_stock):
        self.log(f"Propagating write operation for book {book_id}: {book_name}, quantity: {book_stock}")

        for rep_rank in self.replicas_id:
            if rep_rank == self.rank:  # Skip self
                continue

            try:
                stub = self.get_stub(rep_rank)
                if stub:
                    endpoint = self.replica_endpoints.get(rep_rank)
                    channel = grpc.insecure_channel(endpoint)
                    db_stub = database_grpc.DatabaseServiceStub(channel)

                    response = db_stub.replicateWrite(
                        database.ReplicationRequest(
                            bookID=book_id,
                            bookName=book_name,
                            bookStock=book_stock,
                            sourceReplicaID=self.rank
                        ),
                        timeout=3
                    )

                    self.log(f"Write propagated to replica {rep_rank}: {response.isValid}")
            except Exception as e:
                self.log(f"Failed to propagate write to replica {rep_rank}: {str(e)}")

    def replicateWrite(self, request, context):
        source_replica = request.sourceReplicaID

        with self.lock:
            if source_replica != self.leader_rank:
                self.log(f"Rejecting replication from non-leader replica {source_replica}")
                return database.ReplicationResponse(
                    isValid=False,
                    errMessage="Replication rejected: sender is not the leader"
                )

            self.log(f"Applying replicated write for book {request.bookID}: {request.bookName}")
            self.stock[request.bookID] = [request.bookName, request.bookStock]

            return database.ReplicationResponse(
                isValid=True
            )

    def getStatus(self, request, context):
        with self.lock:
            is_leader = (self.leader_rank == self.rank)
            self.log(f"Status request received. I am {'' if is_leader else 'not '}the leader")

            return replica.StatusResponse(
                replicaID=self.rank,
                isLeader=is_leader,
                leader=self.leader_rank,
                lastAccess=self.last_access,
                knownReplicas=list(self.replicas_id)
            )


def serve():
    service = DatabaseService()
    port = 50060 + service.rank
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    database_grpc.add_DatabaseServiceServicer_to_server(service, server)
    replica_grpc.add_ReplicaServicer_to_server(service, server)
    server.add_insecure_port(f"[::]:{port}")

    service.log(f"Replica started on port {port}")
    server.start()

    server.wait_for_termination()


if __name__ == '__main__':
    serve()