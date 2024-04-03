# gpt-cache : A zero setup semantic cache proxy server for LLM Queries

A work in progress caching solution for saving tokens on LLM queries AND improving responsiveness.

## Benchmark Examples

It's all good and well saying its more response but how much more responsive is it?

See [examples](https://github.com/jwtly10/gpt-cache/tree/main/benchmark-examples/example-go) for real examples of this proxy in action, across various languages and SDKs.

Here's a sneak peak:

Some details are abstracted out for simplicities sake but the full example below can be found [here](www.google.com)

```go
func main() {
    // Using the go-to https://github.com/sashabaranov/go-openai SDK for OpenAI
	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
    // Using the proxy as the base url instead of https://api.openai.com/v1
	config.BaseURL = "http://localhost:8080"
	c := openai.NewClientWithConfig(config)

	var benchmarks []time.Duration

	log.Println("Making the first request (No Cache hit) ...")
	message := "How do you write to file in Java"
    // Make request uses the OpenAI SDK to make a request to the proxy
    // and returns the time it took to make the request
	res1, _ := makeRequest(c, message)
	benchmarks = append(benchmarks, res1)

	log.Println("Making the second request (Cache hit) ...")
	message = "Show me how to write to file in Java"
    // Make request uses the OpenAI SDK to make a request to the proxy
    // and returns the time it took to make the request
	res2, _ := makeRequest(c, message)
	benchmarks = append(benchmarks, res2)

	fmt.Println()

	log.Println("****************** Benchmark Results ******************")
	for i, benchmark := range benchmarks {
		log.Printf("Request %d took: %s\n", i+1, benchmark)
	}
}
```

Output:

```sh
j@mbp:~/Projects/gpt-cache/examples/example-go$ go run .
2024/04/03 18:03:01 Making the first request (No Cache hit) ...
2024/04/03 18:03:06 Response Length: 1155
2024/04/03 18:03:06 Request took: 5.507892679s

2024/04/03 18:03:06 Making the second request (Cache hit) ...
2024/04/03 18:03:07 Response Length: 1155
2024/04/03 18:03:07 Request took: 135.873053ms

2024/04/03 18:03:07 ****************** Benchmark Results ******************
2024/04/03 18:03:07 Request 1 took: 5.507892679s
2024/04/03 18:03:07 Request 2 took: 135.873053ms
```

As you can see, for a simple exact cache match request, the first request took 5.507892679s and the second request took 135.873053ms.

_Thats 41x times faster!_

<!-- TODO: Create a benchmarking utility where you can potentially run  -->

## How it works

While the proxy server is written in Go, a Python microservice plays a pivotal role in enhancing response times and reducing token usage.

This is achieved by converting text from GPT requests into vector embeddings, which are essentially high-dimensional representations that capture the nuanced semantic meanings of the texts. These embeddings are then indexed using Facebook's FAISS, a performant library designed for fast similarity search and efficient clustering of large datasets.

FAISS excels in identifying the nearest neighbors for a given query in the embedding space, allowing our system to perform semantic similarity searches. This means that when a new GPT request is received, its embedding is compared against the embeddings of previously cached requests. The system then identifies the most semantically similar requests already in the cache. This approach moves beyond simple keyword matching, enabling the identification of relevant responses even when the exact wording of requests differs, thereby leveraging the context and meaning behind user queries.

However, semantic caching, as powerful as it is, introduces the possibility of false positives and false negatives. A false positive occurs when the system incorrectly identifies a cache hit, retrieving a semantically similar but ultimately irrelevant response. In a LLM powered system, this is a terrible user experience as the response is not only incorrect but also potentially misleading.

Conversely, a false negative happens when the system misses an appropriate cache hit, usually due to the semantic search parameters being too strict or the embeddings not capturing the nuances of similarity adequately. This is less detrimental, so focus is placed on minimizing false positives.

This project aims to provide a caching system that is both performant and accurate, with a focus on the latter.
