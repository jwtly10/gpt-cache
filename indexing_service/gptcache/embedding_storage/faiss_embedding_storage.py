import numpy as np
import faiss
from typing import List, Tuple

from gptcache.embedding_storage import BaseEmbeddingStorage


class FaissEmbeddingStorage(BaseEmbeddingStorage):
    def __init__(self, dimension: int):
        self.dimension = dimension
        # Determine the metric
        self.index = faiss.IndexFlatL2(dimension)
        self.id_map = []  # To keep track of IDs

    def add_item(self, item_id: int, vector: List[float]):
        print(f"Adding item {item_id} to the index.")
        # Ensure vector is a numpy array of the correct shape (2D)
        if not isinstance(vector, np.ndarray):
            vector = np.array(vector, dtype="float32")
        if vector.ndim == 1:  # If the vector is 1D, reshape it to 2D
            vector = vector.reshape(1, -1)
        if vector.shape[1] != self.dimension:
            raise ValueError("Vector dimension mismatch.")
        self.index.add(vector)
        self.id_map.append(item_id)

    def get_nns_by_vector(
        self, vector: List[float], n: int = 10
    ) -> Tuple[List[int], List[float]]:
        if len(self.id_map) == 0:
            print("Index is empty.")
            return ([], [])

        if not isinstance(vector, np.ndarray):
            vector = np.array(vector, dtype="float32").reshape(1, -1)
        if vector.ndim == 1:  # If the vector is 1D, reshape it to 2D
            vector = vector.reshape(1, -1)
        distances, indices = self.index.search(vector, n)

        print(f"Found {len(indices[0])} neighbors.")
        print(f"Nearest neighbor IDs: {indices[0]}")
        print(f"Distances: {distances[0]}")
        print(f"ID map: {self.id_map}")

        return ([self.id_map[idx] for idx in indices[0]], distances[0].tolist())

    def build_index(self, num_trees: int):
        # Not applicable for IndexFlatL2, but can be used with other FAISS indexes.
        pass

    def save_index(self, filepath: str):
        faiss.write_index(self.index, filepath)

    def load_index(self, filepath: str):
        self.index = faiss.read_index(filepath)
        # Note: You will need to also save and load self.id_map to fully restore state.
