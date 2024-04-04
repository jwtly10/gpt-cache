# gpt-cache : A zero setup semantic cache proxy server for LLM Queries

🏗️ A work in progress caching solution for saving tokens on LLM queries AND improving responsiveness.

3 main goals:

-   Very minimal setup
-   Accurate
-   Easy to optimize

## Benchmark Examples

How much more responsive is it really?

See [examples](https://github.com/jwtly10/gpt-cache/tree/main/benchmark-examples) for real examples of this proxy in action, across various languages and SDKs.

Here's an example of what this proxy can do:

Some details are abstracted out for simplicities sake but the full code example below can be found [here.](https://github.com/jwtly10/gpt-cache/tree/main/benchmark-examples/example-go)

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
    // makeRequest() uses the OpenAI SDK to make a request to the proxy
    // and returns the time it took to make the request
	res1, _ := makeRequest(c, message)
	benchmarks = append(benchmarks, res1)

	log.Println("Making the second request (Cache hit) ...")
	message = "Show me how to write to file in Java"
    // makeRequest() uses the OpenAI SDK to make a request to the proxy
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
j@mbp:~/Projects/gpt-cache/benchmark-examples/example-go$  go run .
2024/04/04 18:15:43 Making the first request (No Cache hit) ...
2024/04/04 18:15:52 Response Length: 1235
2024/04/04 18:15:52 Request took: 8.681519257s

2024/04/04 18:15:52 Making the second request (Cache hit) ...
2024/04/04 18:15:52 Response Length: 1235
2024/04/04 18:15:52 Request took: 66.482852ms

2024/04/04 18:15:52 ****************** Benchmark Results ******************
2024/04/04 18:15:52 Request 1 took: 8.681519257s
2024/04/04 18:15:52 Request 2 took: 66.482852ms
```

As you can see, for a simple exact cache match request, the first request took 8.681519257s and the second request took 66.482852ms.

_Thats 130x times faster!_

<!-- TODO: Create a benchmarking utility where you can potentially run  -->

## How it works

While the proxy server is written in Go, a Python microservice plays a pivotal role in enhancing response times and reducing token usage.

This is achieved by converting text from GPT requests into vector embeddings, which are essentially high-dimensional representations that capture the nuanced semantic meanings of the texts. These embeddings are then indexed using Facebook's FAISS, a performant library designed for fast similarity search and efficient clustering of large datasets.

FAISS excels in identifying the nearest neighbors for a given query in the embedding space, allowing our system to perform semantic similarity searches. This means that when a new GPT request is received, its embedding is compared against the embeddings of previously idexed requests. The system then identifies the most semantically similar requests already in the index. This approach moves beyond simple keyword matching, enabling the identification of relevant responses even when the exact wording of requests differs, thereby leveraging the context and meaning behind user queries.

However, semantic caching, as powerful as it is, introduces the possibility of false positives and false negatives. A false positive occurs when the system incorrectly identifies a cache hit, retrieving a semantically similar but ultimately irrelevant response. In a LLM powered system, this is a terrible user experience as the response is not only incorrect but also potentially misleading.

Conversely, a false negative happens when the system misses an appropriate cache hit, usually due to the semantic search parameters being too strict or the embeddings not capturing the nuances of similarity adequately. This is less detrimental, so focus is placed on minimizing false positives.

This project aims to provide a caching system that is both performant and accurate, with a focus on the latter.

## Architecture

Due to microservices architecture, there are 2 seperate services that need to be run in order to get the proxy server up and running.

-   Proxy Service

    -   This is the main service that acts as the proxy server for the LLM queries.
    -   It proxies the requests to the OpenAI API or returns the response from the cache if it exists based on the semantic similarity of the query.

-   Indexing Service
    -   This service is responsible for indexing the embeddings of the LLM responses and storing them in a FAISS index.
    -   This service is responsible for finding the most semantically similar responses to a given query.

## Setup

### Docker

For ease of development, the services are run in separate docker containers.

Assuming you have Docker/Docker-Compose installed you can run from project root:

```sh
docker-compose up --build
```

to spin up the services, as well as Redis for caching.

Currently the Python image is 7GB in size due to the FAISS library, NLP models and other dependencies. This will be optimized in the future.

### Manual

To run services manually you will need to have Go, and Python 3.12 installed, as well as a Redis instance running for the proxy caching.

Go Service:

```sh
cd proxy_service &&
go run .
```

Some environment variables can be set to configure the proxy:

-   CACHE_THRESHOLD: The threshold for minimum semantic similarity distance between queries.
    -   Default is 0.2
-   REDIS_HOST_URL: The host URL for the Redis instance.
    -   Default is 127.0.0.1 (localhost for docker Redis)
-   REDIS_PORT: The port for the Redis instance.
    -   Default is 6379
-   REDIS_PASSWORD: The password for the Redis instance.
    -   Default is ""
-   REDIS_DB: The database for the Redis instance.
    -   Default is 0
-   INDEX_HOST_URL: The host URL for the Indexing service.
    -   Default is localhost
-   INDEX_PORT: The port for the Indexing service.
    -   Default is 8000

Python Service:

```sh
cd indexing_service &&
pip install -r requirements.txt
uvicorn app.main:app --reload
```

## Contributing

OOS Contributions are welcome! However the project is still in its early stages so please reach out before considering contributing new code.

## License

gpt-cache is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
