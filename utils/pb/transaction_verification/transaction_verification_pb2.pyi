from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class TransactionVerificationRequest(_message.Message):
    __slots__ = ("json",)
    JSON_FIELD_NUMBER: _ClassVar[int]
    json: str
    def __init__(self, json: _Optional[str] = ...) -> None: ...

class TransactionVerificationResponse(_message.Message):
    __slots__ = ("isValid", "errMessage")
    ISVALID_FIELD_NUMBER: _ClassVar[int]
    ERRMESSAGE_FIELD_NUMBER: _ClassVar[int]
    isValid: bool
    errMessage: str
    def __init__(self, isValid: bool = ..., errMessage: _Optional[str] = ...) -> None: ...
