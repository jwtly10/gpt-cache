from pydantic import BaseModel


class QueryResponse(BaseModel):
    id: int = None
    distance: float = None


class AddResponse(BaseModel):
    status: str
    message: str = None
