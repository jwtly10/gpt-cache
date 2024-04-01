from gptcache.embedding_storage import BaseEmbeddingStorage
from annoy import AnnoyIndex
from typing import List, Tuple

class AnnoyEmbeddingStorage(BaseEmbeddingStorage):
    def __init__(self, dimension: int, metric: str = 'angular'):
        self.index = AnnoyIndex(dimension, metric)
        self.dimension = dimension

    def add_item(self, item_id: int, vector: List[float]):
        self.index.add_item(item_id, vector)

    def get_nns_by_vector(self, vector: List[float], n: int = 10) -> List[Tuple[int, float]]:
        return self.index.get_nns_by_vector(vector, n, include_distances=True)

    def build_index(self, num_trees: int):
        self.index.build(num_trees)

    def save_index(self, filepath: str):
        self.index.save(filepath)

    def load_index(self, filepath: str):
        self.index.load(filepath)
