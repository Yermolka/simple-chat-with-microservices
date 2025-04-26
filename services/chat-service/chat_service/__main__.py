from core.services.auth.auth import AuthService
from core.services.chat.chat import ChatService
from core.repository.repository import Repository
from core.repository.db import DB
from core.http.server import app
import uvicorn


auth_service = AuthService()
chat_service = ChatService()
chat_service.inject(Repository(DB("", 0, "", "", "")))

app.dependency_overrides[ChatService] = lambda: chat_service
app.dependency_overrides[AuthService] = lambda: auth_service

if __name__ == "__main__":
    uvicorn.run("__main__:app", host="0.0.0.0", port=3000, reload=True, reload_dirs=["chat_service"])
