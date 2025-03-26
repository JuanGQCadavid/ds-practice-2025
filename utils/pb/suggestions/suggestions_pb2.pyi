from common import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class BookSuggest(_message.Message):
    __slots__ = ("books",)
    class Book(_message.Message):
        __slots__ = ("bookId", "title", "author")
        BOOKID_FIELD_NUMBER: _ClassVar[int]
        TITLE_FIELD_NUMBER: _ClassVar[int]
        AUTHOR_FIELD_NUMBER: _ClassVar[int]
        bookId: str
        title: str
        author: str
        def __init__(self, bookId: _Optional[str] = ..., title: _Optional[str] = ..., author: _Optional[str] = ...) -> None: ...
    BOOKS_FIELD_NUMBER: _ClassVar[int]
    books: _containers.RepeatedCompositeFieldContainer[BookSuggest.Book]
    def __init__(self, books: _Optional[_Iterable[_Union[BookSuggest.Book, _Mapping]]] = ...) -> None: ...
