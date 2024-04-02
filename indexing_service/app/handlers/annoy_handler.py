from gptcache.embedding import SentenceEmbedding
from gptcache.embedding_storage import AnnoyEmbeddingStorage

from app.models.query_model import Query


class AnnoyHandler:
    def __init__(
        self,
        embedding: SentenceEmbedding,
        storage: AnnoyEmbeddingStorage,
    ):
        self.s = embedding
        self.a = storage

    def handle_add(self, query: Query) -> dict:
        """
        Adds a new query's embedding to the Annoy index using the query's unique ID as the key.

        Parameters:
        - query (Query): An instance of the Query model containing the 'id' and 'context'
        of the query. The 'id' is used as the key in the Annoy index, and the 'context'
        is the textual content from which the embedding is generated.

        Returns:
        - dict: A dictionary indicating the result of the operation. It returns
        {"status": "success"} if the embedding is successfully added to the Annoy index.
        In case of an exception, it returns {"status": "error", "message": str(e)},
        where `e` is the exception raised during the operation.

        Raises:
        - Exception: Captures any exceptions raised during the embedding generation or
        when adding the item to the Annoy index, and returns an error message as part
        of the response.
        """
        try:
            # Create embedding for the query
            query_embedding = self.s.to_embedding(query.context)

            # Add the query to the Annoy index
            self.a.add_item(query.id, query_embedding)
        except Exception as e:
            return {"status": "error", "message": str(e)}

        return {"status": "success"}

    def rebuild_index(self):
        """
        Rebuilds the Annoy index in the background after adding a new item.
        """
        self.a.build_index(num_trees=10)
        # TODO: Persist index?
        # annoy.save('your_index_file.ann')

    def handle_query(self, context: str, distance_threshold=0.2) -> dict:
        """
        Handles a query request by generating an embedding for the given context and querying
        the Annoy index to find the closest neighbor and its distance.

        TODO: Currently does not throw ay errors, but should be updated to handle exceptions.
        However, note that we do not need a fallback flow nessasarily, as this is meant for performance.
        If something goes wrong and we get a cache miss, we can live with it. Howver need to investigate
        error conditions here.

        Parameters:
        - context (str): The textual content for which to find the closest matching item
        in the Annoy index based on embedding similarity.
        - distance_threshold (float, optional): The maximum distance between the query
        embedding and a neighbor's embedding for the neighbor to be considered
        sufficiently close. Defaults to 0.2.

        Returns:
        - dict: A dictionary containing the 'id' of the closest matching item (or `None`
        if no such item is found within the threshold) and the 'distance' to this item
        (or `None` if no item is within the threshold).
        """
        print(f"Querying with distance threshold: {distance_threshold}")

        # Create embedding for the query
        query_embedding = self.s.to_embedding(context)

        nn_ids, distances = self.a.get_nns_by_vector(query_embedding, n=1)
        if len(nn_ids) > 0:
            print(f"Nearest neighbor ID: {nn_ids[0]}")
            print(f"Distance: {distances[0]}")
            # Check if the closest neighbor is within the acceptable distance
            if distances[0] <= distance_threshold:
                print("Nearest neighbor found within the distance threshold.")
                return {"id": nn_ids[0], "distance": distances[0]}
            else:
                print("No neighbors found within the distance threshold.")
                return {"id": None, "distance": None}
        else:
            # Handles cases where the index is empty
            print("No neighbors found in the index.")
            return {"id": None, "distance": None}
