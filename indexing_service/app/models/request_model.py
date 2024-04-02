from pydantic import BaseModel


class AddRequest(BaseModel):
    id: int
    context: str


class QueryRequest(BaseModel):
    context: str
    distance_threshold: float = 0.2
