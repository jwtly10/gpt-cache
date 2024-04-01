from abc import ABC, abstractmethod

class BaseEmbedding(ABC):
    @abstractmethod
    def to_embedding(self, text: str) -> list:
        pass


