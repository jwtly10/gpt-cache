import pytest

from app.models.query_model import Query

from app.handlers import AnnoyHandler

from gptcache.embedding_storage import AnnoyEmbeddingStorage
from gptcache.embedding import SentenceEmbedding


def test_handle_add_integration():
    # Setup
    sentence_embedding = SentenceEmbedding(model_name="all-MiniLM-L6-v2")
    annoy = AnnoyEmbeddingStorage(dimension=384)

    handler = AnnoyHandler(sentence_embedding, annoy)

    query_id = 1
    query_context = "The quick brown fox jumps over the lazy dog."

    test_query = Query(id=query_id, context=query_context)

    # Test
    response = handler.handle_add(test_query)
    if response["status"] == "error":
        print(response["message"])

    assert response == {"status": "success"}


@pytest.fixture
def setup_annoy_index():
    """
    Set up and build an Annoy index with pre-defined test sentences.

    Returns:
        AnnoyEmbeddingStorage: An instance of AnnoyEmbeddingStorage with the built index.
        SentenceEmbedding: An instance of SentenceEmbedding for generating embeddings.
    """
    annoy = AnnoyEmbeddingStorage(dimension=384)
    sentence_embedding = SentenceEmbedding(model_name="all-MiniLM-L6-v2")

    test_sentences = [
        "The quick brown fox jumps over the lazy dog.",
        "A fast, dark-colored fox leaps above a sleepy canine.",
        "Today is a sunny day.",
        "I love reading books about science.",
        "Mathematics is the language of the universe.",
        "The language of the universe is mathematics.",
        "Exploring the depths of the ocean reveals hidden treasures.",
        "Artificial intelligence will shape the future of humanity.",
        "A gentle breeze sways the tall trees in the forest.",
        "Culinary skills are both an art and a science.",
    ]

    for i, sentence in enumerate(test_sentences, start=1):
        embedding = sentence_embedding.to_embedding(sentence)
        annoy.add_item(i, embedding)

    annoy.build_index(num_trees=10)

    return annoy, sentence_embedding


def test_query_exact_match(setup_annoy_index):
    annoy, embedding = setup_annoy_index

    handler = AnnoyHandler(embedding, annoy)

    context = "The language of the universe is mathematics."
    distance_threshold = 0.0

    # returns {id: int/None, distance:float/None }
    response = handler.handle_query(context, distance_threshold)

    assert response["id"] == 6
    assert response["distance"] < 0.01
