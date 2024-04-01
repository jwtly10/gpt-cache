from fastapi import FastAPI
from pydantic import BaseModel

from gptcache.embedding import SentenceEmbedding
from gptcache.embedding_storage import AnnoyEmbeddingStorage

app = FastAPI()

s = SentenceEmbedding(model_name="all-MiniLM-L6-v2")
annoy = AnnoyEmbeddingStorage(dimension=384)


class Query(BaseModel):
    context: str
    id: str


@app.post("/query")
async def query(query: Query):
    """
    Queries the closest matching item within the Annoy index based on the embedding similarity of the provided text context.

    This endpoint accepts a request containing a 'context' field with text content. It generates an embedding for the given text and queries the Annoy index to find the closest neighbor based on this embedding. The similarity between the query embedding and the closest neighbor is assessed using a predefined distance threshold.

    If the closest neighbor's distance to the query embedding is within the threshold, the endpoint returns the ID of the neighbor and the distance. If no neighbor is within the threshold, it returns `None` for both the ID and the distance.

    Parameters:
    - context (str): The textual content for which a similar item is being queried.

    Returns:
    - A JSON object containing:
      - 'id': The ID of the closest matching item if one is found within the threshold; otherwise, `None`.
      - 'distance': The distance of the closest matching item to the query; `None` if no item is within the threshold.
    """

    # Create embedding for the query
    query_embedding = s.to_embedding(query.context)

    # Get the closest neighbor and its distance
    nn_id, distance = annoy.get_nns_by_vector(
        query_embedding, n=1, include_distances=True
    )

    # Similarity threshold
    distance_threshold = 0.2

    # Check if the closest neighbor is within the acceptable distance
    if distance[0] <= distance_threshold:
        return {"id": nn_id[0], "distance": distance[0]}
    else:
        return {"id": None, "distance": None}


def rebuild_index():
    """
    Rebuilds the Annoy index in the background after adding a new item.
    """
    annoy.build_index(num_trees=10)
    # TODO: Persist index?
    # annoy.save('your_index_file.ann')


@app.post("/add")
async def add(query: Query, background_task: rebuild_index):
    """
    Adds a new item to the Annoy index based on the provided text context and schedules an index rebuild.

    This endpoint processes a request with a 'context' field, generates an embedding for the text, and adds this embedding to the Annoy index using a unique ID. After adding the item, it schedules an asynchronous task to rebuild the index, ensuring that the new item can be queried.

    Note: The actual rebuilding of the index is handled as a background task to minimize response time and improve the user experience. Considerations regarding when to persist the rebuilt index to disk should be addressed to ensure data durability across application restarts.

    Parameters:
    - context (str): The textual content based on which an embedding will be generated and added to the index.
    - id (int): A unique identifier for the new item being added.

    Returns:
    - A JSON object containing 'status': 'success' if the item is successfully added, or 'status': 'error' with an error message if the process fails.
    """

    try:
        # Create embedding for the query
        query_embedding = s.to_embedding(query.context)

        # Add the query to the Annoy index
        annoy.add_item(query.id, query_embedding)
    except Exception as e:
        return {"status": "error", "message": str(e)}

    return {"status": "success"}
