import grpc
from .auth_pb2_grpc import AuthServiceStub
from .auth_pb2 import VerifyTokenRequest, VerifyTokenResponse


class AuthService:
    def __init__(self, host: str = "localhost", port: int = 50051):
        self.user_service_url = f"{host}:{port}"

    def verify_token(self, user_id: int, token: str) -> bool:
        channel = grpc.insecure_channel(self.user_service_url)
        stub = AuthServiceStub(channel)
        request = VerifyTokenRequest(user_id=str(user_id), token=token)
        response: VerifyTokenResponse = stub.VerifyToken(request)
        return response.valid
