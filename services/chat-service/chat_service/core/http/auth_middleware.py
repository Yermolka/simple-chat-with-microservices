from fastapi.applications import FastAPI
from fastapi.exceptions import HTTPException
from starlette.types import ASGIApp
from core.services.auth.auth import AuthService


class AuthMiddleware:
    def __init__(self, app: FastAPI):
        self.app = app

    def __call__(self, app: ASGIApp, /, *args, **kwargs) -> ASGIApp:
        self.auth_service = self.app.dependency_overrides[AuthService]()

        async def middleware(scope, receive, send):
            if not scope.get("headers"):
                raise HTTPException(status_code=401, detail="Unauthorized")
            
            headers = dict(scope["headers"])
            auth_header = headers.get(b"authorization")
            if not auth_header:
                raise HTTPException(status_code=401, detail="Unauthorized")
            
            try:
                auth_type, token = auth_header.decode("utf-8").split(" ")
            except ValueError:
                raise HTTPException(status_code=401, detail="Unauthorized")
            
            if not token or auth_type.lower() != "bearer":
                raise HTTPException(status_code=401, detail="Unauthorized")

            user_id = headers.get(b"x-user-id")
            if not user_id:
                raise HTTPException(status_code=401, detail="Unauthorized")

            if not self.auth_service.verify_token(int(user_id.decode("utf-8")), token):
                raise HTTPException(status_code=401, detail="Unauthorized")

            return await app(scope, receive, send)
        
        return middleware
