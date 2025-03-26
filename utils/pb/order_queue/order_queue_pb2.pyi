from common import common_pb2 as _common_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EmptyRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class EnqueueRequest(_message.Message):
    __slots__ = ("orderId", "order")
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    ORDER_FIELD_NUMBER: _ClassVar[int]
    orderId: str
    order: _common_pb2.Order
    def __init__(self, orderId: _Optional[str] = ..., order: _Optional[_Union[_common_pb2.Order, _Mapping]] = ...) -> None: ...

class EnqueueResponse(_message.Message):
    __slots__ = ("errMessage", "isValid")
    ERRMESSAGE_FIELD_NUMBER: _ClassVar[int]
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    errMessage: str
    isValid: bool
    def __init__(self, errMessage: _Optional[str] = ..., isValid: bool = ...) -> None: ...

class DequeueResponse(_message.Message):
    __slots__ = ("orderId", "order", "errMessage", "isValid")
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    ORDER_FIELD_NUMBER: _ClassVar[int]
    ERRMESSAGE_FIELD_NUMBER: _ClassVar[int]
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    orderId: str
    order: _common_pb2.Order
    errMessage: str
    isValid: bool
    def __init__(self, orderId: _Optional[str] = ..., order: _Optional[_Union[_common_pb2.Order, _Mapping]] = ..., errMessage: _Optional[str] = ..., isValid: bool = ...) -> None: ...
