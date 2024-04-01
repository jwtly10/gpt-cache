from fastapi import FastAPI
from pydantic import BaseModel

from nlp import to_similarity

app = FastAPI()

@app.get("/")
async def root():
    return "Healthy"

class Compare(BaseModel):
    a: str
    b: str

@app.post("/getSimilarity")
async def get_similarity(comp: Compare):
    return {"similarity": to_similarity(comp.a, comp.b)}
