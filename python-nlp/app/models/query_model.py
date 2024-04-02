from pydantic import BaseModel


class Query(BaseModel):
    context: str
    id: int
