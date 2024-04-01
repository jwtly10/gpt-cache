from gptcache.embedding import BaseEmbedding
from sentence_transformers import SentenceTransformer

class SentenceEmbedding(BaseEmbedding):
    def __init__(self, model_name: str):
        self.model = SentenceTransformer(model_name)

    def to_embedding(self, text: str):
        return self.model.encode(text)