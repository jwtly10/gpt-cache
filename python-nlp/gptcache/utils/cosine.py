import torch

def cosine_similarity(embedding1, embedding2):
    """
    Compute the cosine similarity between two embeddings using PyTorch.

    Args:
    - embedding1 (torch.Tensor or np.ndarray or list): First embedding vector.
    - embedding2 (torch.Tensor or np.ndarray or list): Second embedding vector.

    Returns:
    - float: Cosine similarity score.
    """
    # Ensure the embeddings are PyTorch tensors
    if not isinstance(embedding1, torch.Tensor):
        embedding1 = torch.tensor(embedding1)
    if not isinstance(embedding2, torch.Tensor):
        embedding2 = torch.tensor(embedding2)
    
    # Normalize the embeddings to unit vectors
    embedding1_norm = embedding1 / embedding1.norm()
    embedding2_norm = embedding2 / embedding2.norm()
    
    # Compute the cosine similarity
    similarity = torch.dot(embedding1_norm, embedding2_norm).item()  # Convert to Python float
    
    return similarity
