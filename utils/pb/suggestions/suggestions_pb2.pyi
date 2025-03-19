from common import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ItemsBought(_message.Message):
    __slots__ = ("items",)
    class Item(_message.Message):
        __slots__ = ("name", "quantity")
        NAME_FIELD_NUMBER: _ClassVar[int]
        QUANTITY_FIELD_NUMBER: _ClassVar[int]
        name: str
        quantity: str
        def __init__(self, name: _Optional[str] = ..., quantity: _Optional[str] = ...) -> None: ...
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[ItemsBought.Item]
    def __init__(self, items: _Optional[_Iterable[_Union[ItemsBought.Item, _Mapping]]] = ...) -> None: ...

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
