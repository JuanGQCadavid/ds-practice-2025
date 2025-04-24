from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class InitRequest(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class InitResponse(_message.Message):
    __slots__ = ("message", "isValid")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    message: str
    isValid: bool
    def __init__(self, message: _Optional[str] = ..., isValid: bool = ...) -> None: ...

class StockRequest(_message.Message):
    __slots__ = ("bookID", "bookName", "bookStock")
    BOOKID_FIELD_NUMBER: _ClassVar[int]
    BOOKNAME_FIELD_NUMBER: _ClassVar[int]
    BOOKSTOCK_FIELD_NUMBER: _ClassVar[int]
    bookID: str
    bookName: str
    bookStock: int
    def __init__(self, bookID: _Optional[str] = ..., bookName: _Optional[str] = ..., bookStock: _Optional[int] = ...) -> None: ...

class StockResponse(_message.Message):
    __slots__ = ("bookName", "bookStock", "isValid", "errMessage", "leaderID")
    BOOKNAME_FIELD_NUMBER: _ClassVar[int]
    BOOKSTOCK_FIELD_NUMBER: _ClassVar[int]
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    ERRMESSAGE_FIELD_NUMBER: _ClassVar[int]
    LEADERID_FIELD_NUMBER: _ClassVar[int]
    bookName: str
    bookStock: int
    isValid: bool
    errMessage: str
    leaderID: int
    def __init__(self, bookName: _Optional[str] = ..., bookStock: _Optional[int] = ..., isValid: bool = ..., errMessage: _Optional[str] = ..., leaderID: _Optional[int] = ...) -> None: ...

class ReplicationRequest(_message.Message):
    __slots__ = ("bookID", "bookName", "bookStock", "sourceReplicaID")
    BOOKID_FIELD_NUMBER: _ClassVar[int]
    BOOKNAME_FIELD_NUMBER: _ClassVar[int]
    BOOKSTOCK_FIELD_NUMBER: _ClassVar[int]
    SOURCEREPLICAID_FIELD_NUMBER: _ClassVar[int]
    bookID: str
    bookName: str
    bookStock: int
    sourceReplicaID: int
    def __init__(self, bookID: _Optional[str] = ..., bookName: _Optional[str] = ..., bookStock: _Optional[int] = ..., sourceReplicaID: _Optional[int] = ...) -> None: ...

class ReplicationResponse(_message.Message):
    __slots__ = ("isValid", "errMessage")
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    ERRMESSAGE_FIELD_NUMBER: _ClassVar[int]
    isValid: bool
    errMessage: str
    def __init__(self, isValid: bool = ..., errMessage: _Optional[str] = ...) -> None: ...

class StatusResponse(_message.Message):
    __slots__ = ("isLeader", "replicaID", "leaderID")
    ISLEADER_FIELD_NUMBER: _ClassVar[int]
    REPLICAID_FIELD_NUMBER: _ClassVar[int]
    LEADERID_FIELD_NUMBER: _ClassVar[int]
    isLeader: bool
    replicaID: int
    leaderID: int
    def __init__(self, isLeader: bool = ..., replicaID: _Optional[int] = ..., leaderID: _Optional[int] = ...) -> None: ...

class EmptyRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
