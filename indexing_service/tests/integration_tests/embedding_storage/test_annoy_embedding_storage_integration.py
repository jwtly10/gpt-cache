import pytest

from gptcache.embedding_storage import AnnoyEmbeddingStorage
from gptcache.embedding import SentenceEmbedding


@pytest.fixture
def annoy_storage_and_embeddings():
    # Initialize the sentence embedding model
    sentence_embedding = SentenceEmbedding(model_name="all-MiniLM-L6-v2")

    # Initialize AnnoyEmbeddingStorage
    annoy_storage = AnnoyEmbeddingStorage(
        dimension=384
    )  # Dimension matches the model's output

    # Sentences setup
    closely_related_sentence1 = "The quick brown fox jumps over the lazy dog."
    closely_related_sentence2 = "A fast, dark-colored fox leaps above a sleepy canine."
    unrelated_sentences = [
        "Today is a sunny day.",
        "I love reading books about science.",
        "Mathematics is the language of the universe.",
    ]

    # Generate embeddings for all sentences and add them to Annoy
    sentences = [
        closely_related_sentence1,
        closely_related_sentence2,
    ] + unrelated_sentences
    for i, sentence in enumerate(sentences, start=1):
        embedding = sentence_embedding.to_embedding(sentence)
        annoy_storage.add_item(i, embedding)

    # Build the index
    annoy_storage.build_index(num_trees=10)

    return annoy_storage, sentence_embedding


def test_find_closely_related_sentence_integration(annoy_storage_and_embeddings):
    annoy_storage, sentence_embedding = annoy_storage_and_embeddings
    query_embedding = sentence_embedding.to_embedding(
        "The quick fox who is brown jumps over the dog who is lazy"
    )
    nearest_neighbors = annoy_storage.get_nns_by_vector(query_embedding, n=1)
    nearest_neighbor_id = nearest_neighbors[0][0]
    assert (
        nearest_neighbor_id == 1
    ), f"Expected the nearest neighbor to be ID 1, got {nearest_neighbor_id}."

    annoy_storage, sentence_embedding = annoy_storage_and_embeddings
    query_embedding = sentence_embedding.to_embedding(
        "The language of the universe is mathematics."
    )
    nearest_neighbors = annoy_storage.get_nns_by_vector(query_embedding, n=1)
    nearest_neighbor_id = nearest_neighbors[0][0]
    assert (
        nearest_neighbor_id == 5
    ), f"Expected the nearest neighbor to be ID 5, got {nearest_neighbor_id}."
