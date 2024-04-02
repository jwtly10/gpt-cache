import pytest

from gptcache.utils.cosine import cosine_similarity

from gptcache.embedding.sentence_embedding import SentenceEmbedding


def test_get_similarity_of_embeddings():
    sentence1 = "Write a function in Python that writes some text to a file."
    sentence2 = "Show me how to write text to a file using Python."

    s = SentenceEmbedding(model_name="all-MiniLM-L6-v2")

    embedding1 = s.to_embedding(sentence1)
    embedding2 = s.to_embedding(sentence2)

    similarity = cosine_similarity(embedding1, embedding2)

    assert (
        similarity > 0.8
    ), f"Expected similarity to be greater than 0.8, got {similarity}."
