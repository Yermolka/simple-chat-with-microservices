from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class VerifyTokenRequest(_message.Message):
    __slots__ = ("user_id", "token")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    token: str
    def __init__(self, user_id: _Optional[str] = ..., token: _Optional[str] = ...) -> None: ...

class VerifyTokenResponse(_message.Message):
    __slots__ = ("valid", "error")
    VALID_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    valid: bool
    error: str
    def __init__(self, valid: bool = ..., error: _Optional[str] = ...) -> None: ...
