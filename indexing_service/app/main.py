from fastapi import FastAPI, BackgroundTasks, HTTPException

from gptcache.embedding import SentenceEmbedding
from gptcache.embedding_storage import FaissEmbeddingStorage

from app.handlers.faiss_handler import FaissHandler
from app.models.request_model import AddRequest, QueryRequest
from app.models.response_model import AddResponse, QueryResponse


app = FastAPI()

embedding = SentenceEmbedding(model_name="all-MiniLM-L6-v2")
faiss = FaissEmbeddingStorage(dimension=384)

handler = FaissHandler(embedding, faiss)


@app.post("/queryIndex", response_model=QueryResponse)
async def query_index(query: QueryRequest):
    res = handler.handle_query(query.context, query.distance_threshold)
    if res["id"] is not None:
        return QueryResponse(id=res["id"], distance=res["distance"])
    else:
        raise HTTPException(status_code=204, detail="No similar context found")


@app.post("/addIndex", response_model=AddResponse)
async def add_index(query: AddRequest):
    try:
        res = handler.handle_add(query.id, query.context)
        if res["status"] == "success":
            return AddResponse(status="success")
        else:
            raise HTTPException(status_code=500, detail=res["message"])
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
