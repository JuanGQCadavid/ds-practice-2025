from common import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class FraudDetectionRequestInit(_message.Message):
    __slots__ = ("orderId", "order")
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    ORDER_FIELD_NUMBER: _ClassVar[int]
    orderId: str
    order: _common_pb2.Order
    def __init__(self, orderId: _Optional[str] = ..., order: _Optional[_Union[_common_pb2.Order, _Mapping]] = ...) -> None: ...

class FraudDetectionRequestClock(_message.Message):
    __slots__ = ("orderId", "clock")
    ORDERID_FIELD_NUMBER: _ClassVar[int]
    CLOCK_FIELD_NUMBER: _ClassVar[int]
    orderId: str
    clock: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, orderId: _Optional[str] = ..., clock: _Optional[_Iterable[int]] = ...) -> None: ...

class FraudDetectionResponse(_message.Message):
    __slots__ = ("code",)
    CODE_FIELD_NUMBER: _ClassVar[int]
    code: str
    def __init__(self, code: _Optional[str] = ...) -> None: ...

class FraudDetectionResponseClock(_message.Message):
    __slots__ = ("response", "clock")
    RESPONSE_FIELD_NUMBER: _ClassVar[int]
    CLOCK_FIELD_NUMBER: _ClassVar[int]
    response: FraudDetectionResponse
    clock: _containers.RepeatedScalarFieldContainer[int]
    def __init__(self, response: _Optional[_Union[FraudDetectionResponse, _Mapping]] = ..., clock: _Optional[_Iterable[int]] = ...) -> None: ...
