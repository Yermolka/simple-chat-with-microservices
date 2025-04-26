from typing import Annotated
from fastapi import FastAPI, Depends
from core.http.auth_middleware import AuthMiddleware
from core.services.chat.chat import ChatService

app = FastAPI()

app.add_middleware(AuthMiddleware(app))

@app.get("/api/chats/", response_model=None)
async def get_chats(chat_service: Annotated[ChatService, Depends()]) -> list:
    return await chat_service.get_chats(0)
