from transformers import AutoTokenizer, AutoModel
import torch
from scipy.spatial.distance import cosine


def to_similarity(sentence1: str, sentence2: str):
    print(f"Comparing: {sentence1} and {sentence2}")

    embedding1 = encode(sentence1)
    embedding2 = encode(sentence2)

    similarity = 1 - cosine(embedding1.detach().numpy(), embedding2.detach().numpy())
    print(f"Similarity: {similarity}")
    return similarity


def encode(sentence):
    tokenizer = AutoTokenizer.from_pretrained("sentence-transformers/all-MiniLM-L6-v2")
    model = AutoModel.from_pretrained("sentence-transformers/all-MiniLM-L6-v2")

    inputs = tokenizer(sentence, return_tensors="pt", padding=True, truncation=True, max_length=512)
    outputs = model(**inputs)
    return outputs.last_hidden_state.mean(dim=1).squeeze()
