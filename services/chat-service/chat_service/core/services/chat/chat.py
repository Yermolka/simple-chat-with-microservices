from core.repository.repository import Repository


class ChatService:
    repository: Repository

    def __init__(self):
        pass

    def inject(self, repository: Repository):
        self.repository = repository
        
    async def get_chats(self, user_id: int) -> list:
        return await self.repository.get_chats(user_id)
