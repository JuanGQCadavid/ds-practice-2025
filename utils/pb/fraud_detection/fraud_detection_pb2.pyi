from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class CreditCard(_message.Message):
    __slots__ = ("number", "cvv", "expirationDate")
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    CVV_FIELD_NUMBER: _ClassVar[int]
    EXPIRATIONDATE_FIELD_NUMBER: _ClassVar[int]
    number: str
    cvv: str
    expirationDate: str
    def __init__(self, number: _Optional[str] = ..., cvv: _Optional[str] = ..., expirationDate: _Optional[str] = ...) -> None: ...

class FraudDetectionRequest(_message.Message):
    __slots__ = ("creditCard",)
    CREDITCARD_FIELD_NUMBER: _ClassVar[int]
    creditCard: CreditCard
    def __init__(self, creditCard: _Optional[_Union[CreditCard, _Mapping]] = ...) -> None: ...

class FraudDetectionResponse(_message.Message):
    __slots__ = ("code",)
    CODE_FIELD_NUMBER: _ClassVar[int]
    code: str
    def __init__(self, code: _Optional[str] = ...) -> None: ...
