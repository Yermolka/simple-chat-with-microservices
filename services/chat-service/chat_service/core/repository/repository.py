from typing import Any
from core.repository.db import DB
from psycopg.cursor_async import AsyncCursor


class Repository:
    db: DB

    def __init__(self, db: DB):
        self.db = db

    async def get_chats(self, user_id: int) -> list:
        return []
    
        # async with self.db.connection() as conn:
        #     async with conn.cursor() as cursor:
        #         return await self._get_chats(cursor, user_id)
    
    async def _get_chats(self, cursor: AsyncCursor[Any], user_id: int) -> list:
        return []
