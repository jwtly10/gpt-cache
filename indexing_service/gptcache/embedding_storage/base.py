from abc import ABC, abstractmethod
from typing import List, Tuple

class BaseEmbeddingStorage(ABC):
    @abstractmethod
    def add_item(self, item_id: int, vector: List[float]):
        pass

    @abstractmethod
    def get_nns_by_vector(self, vector: List[float], n: int = 10) -> List[Tuple[int, float]]:
        pass

    @abstractmethod
    def build_index(self, num_trees: int):
        pass

    @abstractmethod
    def save_index(self, filepath: str):
        pass

    @abstractmethod
    def load_index(self, filepath: str):
        pass
