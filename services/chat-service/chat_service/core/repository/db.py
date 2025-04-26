from psycopg_pool import AsyncConnectionPool
from psycopg.rows import dict_row
from contextlib import asynccontextmanager


class DB:
    pool: AsyncConnectionPool
    
    def __init__(self, host: str, port: int, user: str, password: str, dbname: str):
        pass

        # conn_str = f"postgresql://{user}:{password}@{host}:{port}/{dbname}"
        # self.pool = AsyncConnectionPool(
        #     conn_str,
        #     kwargs={"row_factory": dict_row})
        
    @asynccontextmanager
    async def connection(self, *args, **kwargs):
        async with self.pool.connection(*args, **kwargs) as conn:
            conn.prepare_threshold = 0
            if conn.autocommit is False:
                await conn.set_autocommit(True)

            yield conn
