from common import common_pb2 as _common_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PrepareRequest(_message.Message):
    __slots__ = ("orderID", "creditCard")
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    CREDITCARD_FIELD_NUMBER: _ClassVar[int]
    orderID: str
    creditCard: _common_pb2.CreditCard
    def __init__(self, orderID: _Optional[str] = ..., creditCard: _Optional[_Union[_common_pb2.CreditCard, _Mapping]] = ...) -> None: ...

class CommitRequest(_message.Message):
    __slots__ = ("orderID",)
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    orderID: str
    def __init__(self, orderID: _Optional[str] = ...) -> None: ...

class AbortRequest(_message.Message):
    __slots__ = ("orderID",)
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    orderID: str
    def __init__(self, orderID: _Optional[str] = ...) -> None: ...
