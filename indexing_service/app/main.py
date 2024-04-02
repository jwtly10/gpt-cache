from fastapi import FastAPI

from handlers import add_handler
from handlers import query_handler

from models import Query

app = FastAPI()


@app.post("/query")
async def query(query: Query):
    return query_handler.handle_query(query.context)


@app.post("/add")
async def add(query: Query, background_task: add_handler.rebuild_index):
    return add_handler.handle_add(query)
