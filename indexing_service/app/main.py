from fastapi import FastAPI, BackgroundTasks, HTTPException

from gptcache.embedding import SentenceEmbedding
from gptcache.embedding_storage import AnnoyEmbeddingStorage

from app.handlers.annoy_handler import AnnoyHandler
from app.models.request_model import AddRequest, QueryRequest
from app.models.response_model import AddResponse, QueryResponse


app = FastAPI()

embedding = SentenceEmbedding(model_name="all-MiniLM-L6-v2")
annoy = AnnoyEmbeddingStorage(dimension=384)

handler = AnnoyHandler(embedding, annoy)


# TODO Refactor
@app.post("/queryIndex", response_model=QueryResponse)
async def query_index(query: QueryRequest):
    res = handler.handle_query(query.context, query.distance_threshold)
    if res["id"] is not None:
        return QueryResponse(id=res["id"], distance=res["distance"])
    else:
        return QueryResponse()


@app.post("/addIndex", response_model=AddResponse)
async def add_index(query: AddRequest, background_tasks: BackgroundTasks):
    res = handler.handle_add(query.id, query.context)
    background_tasks.add_task(handler.rebuild_index)

    if res["status"] == "success":
        return AddResponse(status="success")
    else:
        return AddResponse(status="error", message=res["message"])
