from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class IDRequest(_message.Message):
    __slots__ = ("ID",)
    ID_FIELD_NUMBER: _ClassVar[int]
    ID: int
    def __init__(self, ID: _Optional[int] = ...) -> None: ...

class IDResponse(_message.Message):
    __slots__ = ("isValid", "knownIDs")
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    KNOWNIDS_FIELD_NUMBER: _ClassVar[int]
    isValid: bool
    knownIDs: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, isValid: bool = ..., knownIDs: _Optional[_Iterable[int]] = ...) -> None: ...

class LeaderNotification(_message.Message):
    __slots__ = ("leaderID",)
    LEADERID_FIELD_NUMBER: _ClassVar[int]
    leaderID: int
    def __init__(self, leaderID: _Optional[int] = ...) -> None: ...

class LeaderTransfer(_message.Message):
    __slots__ = ("oldLeaderID", "newLeaderID")
    OLDLEADERID_FIELD_NUMBER: _ClassVar[int]
    NEWLEADERID_FIELD_NUMBER: _ClassVar[int]
    oldLeaderID: int
    newLeaderID: int
    def __init__(self, oldLeaderID: _Optional[int] = ..., newLeaderID: _Optional[int] = ...) -> None: ...

class HeartbeatRequest(_message.Message):
    __slots__ = ("leaderID",)
    LEADERID_FIELD_NUMBER: _ClassVar[int]
    leaderID: int
    def __init__(self, leaderID: _Optional[int] = ...) -> None: ...

class HeartbeatResponse(_message.Message):
    __slots__ = ("replicaID", "isAlive")
    REPLICAID_FIELD_NUMBER: _ClassVar[int]
    ISALIVE_FIELD_NUMBER: _ClassVar[int]
    replicaID: int
    isAlive: bool
    def __init__(self, replicaID: _Optional[int] = ..., isAlive: bool = ...) -> None: ...

class StatusRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class StatusResponse(_message.Message):
    __slots__ = ("replicaID", "isLeader", "leader", "lastAccess", "knownReplicas")
    REPLICAID_FIELD_NUMBER: _ClassVar[int]
    ISLEADER_FIELD_NUMBER: _ClassVar[int]
    LEADER_FIELD_NUMBER: _ClassVar[int]
    LASTACCESS_FIELD_NUMBER: _ClassVar[int]
    KNOWNREPLICAS_FIELD_NUMBER: _ClassVar[int]
    replicaID: int
    isLeader: bool
    leader: int
    lastAccess: str
    knownReplicas: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, replicaID: _Optional[int] = ..., isLeader: bool = ..., leader: _Optional[int] = ..., lastAccess: _Optional[str] = ..., knownReplicas: _Optional[_Iterable[int]] = ...) -> None: ...
