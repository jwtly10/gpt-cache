import pytest
from unittest.mock import MagicMock, patch
from gptcache.embedding_storage import (
    AnnoyEmbeddingStorage,
)


@pytest.fixture
def mock_annoy_index():
    with patch("annoy.AnnoyIndex") as mock:
        yield mock()


@pytest.fixture
def embedding_storage(mock_annoy_index):
    storage = AnnoyEmbeddingStorage(dimension=128, metric="angular")
    storage.index = mock_annoy_index
    return storage


def test_add_item(embedding_storage, mock_annoy_index):
    item_id = 1
    vector = [0.1, 0.2, 0.3]
    embedding_storage.add_item(item_id, vector)
    mock_annoy_index.add_item.assert_called_once_with(item_id, vector)


def test_get_nns_by_vector(embedding_storage, mock_annoy_index):
    vector = [0.1, 0.2, 0.3]
    n = 5
    mock_annoy_index.get_nns_by_vector.return_value = ([2, 3], [0.4, 0.5])
    ids, distances = embedding_storage.get_nns_by_vector(vector, n)
    mock_annoy_index.get_nns_by_vector.assert_called_once_with(
        vector, n, include_distances=True
    )
    assert ids == [2, 3]
    assert distances == [0.4, 0.5]


def test_build_index(embedding_storage, mock_annoy_index):
    num_trees = 10
    embedding_storage.build_index(num_trees)
    mock_annoy_index.build.assert_called_once_with(num_trees)


def test_save_index(embedding_storage, mock_annoy_index):
    filepath = "test.ann"
    embedding_storage.save_index(filepath)
    mock_annoy_index.save.assert_called_once_with(filepath)


def test_load_index(embedding_storage, mock_annoy_index):
    filepath = "test.ann"
    embedding_storage.load_index(filepath)
    mock_annoy_index.load.assert_called_once_with(filepath)
