# Twitter-like Service - System Design
This document outlines the system design for a scalable, read-optimized service that mimics core features of Twitter (now X).

This is drafted following the format from a [certification](https://www.educative.io/verify-certificate/lOn30BIPYEZpJv0BrCNBx00MoOqvuM) by [course "grokking the system design interview"](https://www.educative.io/courses/grokking-the-system-design-interview) I completed in January 2025.

## 0. Run API

### Prerequisites
- Docker
- Go

### How to test ?
Run application
1. `docker compose -u`
2. `go mod tidy`
3. Execute main on IDE (todo: build image and add to docker-compose)

### Calling endpoints 

**Follow User**

*Note: X-User-ID header and follow_user_id body params should be exist on the database. Check `init.sql` to find valid IDs*

```
curl --location 'http://localhost:8080/api/v1/follow' \
--header 'X-User-ID: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15' \
--header 'Content-Type: application/json' \
--data '{
    "follow_user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13"
}'
```

**Publish tweet**

*Note: X-User-ID header should be exist on the database. Check `init.sql` to find valid IDs.*

*Note 2: use https://www.uuidgenerator.net/version4 to obtains a idempotency_key valid.*

```
curl --location 'http://localhost:8080/api/v1/tweet' \
--header 'X-User-ID: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11' \
--header 'Content-Type: application/json' \
--data '{
    "text": "test redis #3",
    "idempotency_key": "a00ffe35-fc64-45f3-be60-8c824ec0a346"
}'
```

***Get Timeline***

*Note: X-User-ID header should be exist on the database. Check `init.sql` to find valid IDs.*

```
curl --location 'http://localhost:8080/api/v1/timeline' \
--header 'X-User-ID: a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15'
```

### Test on my laptop
<img width="1321" height="386" alt="image" src="https://github.com/user-attachments/assets/0c57ec3c-21df-4328-a5c7-67523537a3cb" />
<img width="1329" height="805" alt="image" src="https://github.com/user-attachments/assets/86333e69-ece3-4f3d-865f-578c2e0e2842" />
<img width="1320" height="599" alt="image" src="https://github.com/user-attachments/assets/6a56d154-07bb-4c11-b2f1-0cf07c3a7bbc" />




## 1. Business Logic

The core of the system revolves around users following other users. Naturally, some users will have many more followers than they follow. This indicates that the system will be **read-heavy**, with a large number of users consuming tweets and a smaller percentage actively creating them. Tweets are short messages, typically up to 280 characters.

## 2. Requirements

### Functional Requirements

- **Tweeting**: Users can publish short messages (tweets) not exceeding a character limit (e.g., 280 characters).
- **Follow**: Users can follow other users.
- **Timeline**: Users can view a timeline displaying tweets from the users they follow.
- **Authentication Assumption**: All users are considered valid. There is no need to implement a sign-in module or session management. A user's identifier can be passed via a header, parameter, or body as deemed convenient.

### Non-Functional Requirements

- The solution must be able to **scale to millions of users**.
- The application must be **optimized for reads**.

## 3. Scalability Estimates

In 2025, X (formerly Twitter) has approximately 611 million monthly active users and 245 million daily active users.

Source consulted: https://www.demandsage.com/twitter-statistics/

Let's assume an average user:

- Follows ~100 accounts.
- Publishes 10 tweets per day.
- Reads 100 tweets per day.

Read Volume: 245M users * 100 reads/day = 24.5 billion read operations per day.

Write Volume: 245M users * 10 tweets/day = 2.45 billion write operations per day.

### Data Storage Estimates

Assuming each tweet (text-only) averages `1 KB`:

**Write Storage Volume:** 2.45 billion writes/day * 1 KB/write = **2.45 Terabytes (TB) per day**.

### Key Conclusions

- The system is **read-heavy**.
- A significant amount of data needs to be persisted daily.

## 4. High-Level Architecture

// TODO fix

Made by [excalidraw.com](http://excalidraw.com)

## 5. Data Model

A hybrid data model is proposed to leverage the strengths of both relational and NoSQL databases. A relational database (PostgreSQL) will serve as the source of truth for its consistency and ability to handle complex relationships, while a NoSQL database (Redis) will be used as a high-speed cache for read-heavy operations like timeline generation.

My decision to use PostgreSQL stems from my experience at [Sitrack.com](http://sitrack.com/). There, we used it to reliably store millions of records for their vehicle logistics and tracking platform, so I am confident in its ability to scale. I chose Redis based on my experience using it at Mercado Libre and PedidosYa.

### Relational Model (PostgreSQL)

### `Users` Table

- `id` (UUID v4, Primary Key)
- `username` (string)
- `created_at` (timestamp)
- `updated_at` (timestamp)

### `Follows` Table

- `follower_id` (UUID v4, Composite Primary Key, Foreign Key to `Users.id`) - The user who follows.
- `following_id` (UUID v4, Composite Primary Key, Foreign Key to `Users.id`) - The user who is being followed.
- `created_at` (timestamp)

This structure allows for efficient queries, such as `SELECT follower_id FROM Follows WHERE following_id = 'user-x-id';`, to find all followers of a user.


### `Tweets` Table

- `id` (UUID v4, Primary Key)
- `user_id` (UUID v4, Foreign Key to `Users.id`)
- `content` (string)
- `created_at` (timestamp)

### NoSQL Model (Redis)

### User Timeline Cache

- **Key**: `timeline:<user_id>` (e.g., `timeline:f4691a93-f2c0-4480-8172-39f5a9b0105e`)
- **Value**: A Redis List of tweet IDs (e.g., `["tweet_id_34", "tweet_id_12", "tweet_id_99", ...]`)

*Note: A data eviction policy should be defined for these keys to manage memory usage.*

## 6. API Endpoint Design

We assume that all incoming requests contain a user identifier in the `X-User-ID` header.

### Publish a Tweet

- Endpoint `POST /api/v1/tweet`
- Header

```
X-User-ID: "f4691a93-f2c0-4480-8172-39f5a9b0105e"
```

- Request body

```json
{
	"text": "Example tweet", // 280 characters, Required
	"idempontecy_key": "f4691a93-f2c0-4480-8172-39f5a9b0105e" // will be tweet id. Use for avoid duplication and retry for clients
}
```

- Success Response

```json
{
	"tweet_id": "f4691a93-f2c0-4480-8172-39f5a9b0105e",
	"created_at": "2023-09-24T15:30:00Z"
}
```

Response Code Errors

```
201 Created
409 Bad Request
500 Internal Server Error
```

- Validations
    - Text maximum 280 characters
    - the user_id exist
    - the tweet is not already created. Check idempotency_key.

### Follow a User

- Endpoint `POST /api/v1/follow`
- Header

```
X-User-ID: "userID"
```

- Request body

```json
{
	"follow_user_id": "f4691a93-f2c0-4480-8172-39f5a9b0105e" 
}
```

- Response Payload

```json
// only return error code or error message
```

- Response Code Errors

```
201 Created
409 Bad Request
500 Internal Server Error
```

- Validations
    - The user_id and follow_user_id exist
    - Che before if this relation already exist. If this relations already exist, return 201 Created
    - No se puede auto seguir el usuario

### View Timeline

- Endpoint `GET /api/v1/timeline?limit=xx&next_cursor=xxxx`
- Request Header

```jsx
Header
X-User-ID: xxxx
```

- Success Response

```json
{
	"tweets": [
		{
			"tweet_id": "f4691a93-f2c0-4480-8172-39f5a9b0105f",
			"user_id": "f4691a93-f2c0-4480-8172-39f5a9b0105f",
			"text": "Hola soy Elon y ahora se llamará X",
			"created_at": "2023-09-25T15:30:00Z", // Tweet más reciente
			"username": "elonaitor"
		},
		{
			"tweet_id": "f4691a93-f2c0-4480-8172-39f5a9b0105e",
			"user_id": "f4691a93-f2c0-4480-8172-39f5a9b0105e",
			"text": "primer tweet de la historia",
			"created_at": "2023-09-24T15:30:00Z", // Tweet más viejo
			"username": "renzonaitor"
		},
	],
	"next_cursor": "f4691a93-f2c0-4480-8172-39f5a9b0105f" // Tweet próximo
}
```

- Response Code Errors

```
200 OK
409 Bad Request
500 Internal Server Error
```

- Validations
    - Check if the `user_id` exist

## 7. Timeline Generation Flow: "Fan-out on Write"

To ensure the system is highly optimized for reads, we use a **"Fan-out on Write"** (or Push) model.

1. **A user publishes a tweet**:
    - The tweet is saved to the `Tweets` table in the relational database (PostgreSQL).
    - An asynchronous event is dispatched (e.g., via a message queue or a goroutine) containing the `tweet_id` and the author's `user_id`.
2. **A Timeline Worker processes the event**:
    - The worker consumes the event.
    - It queries the `Follows` table in the database to get a list of all `follower_id` for the author.
    - For each `follower_id`, the worker executes the `LPUSH` command in Redis, pushing the new `tweet_id` onto the top of that follower's timeline list.
3. **The `GET /timeline` endpoint becomes extremely performant**:
    - It fetches a list of `tweet_id` from Redis using `LRANGE`. Cursor-based pagination is used to get the correct slice of the list.
    - It "hydrates" these IDs by fetching the full tweet objects from PostgreSQL with a single `SELECT * FROM Tweets WHERE id IN (...)` query. This query is very fast as it uses the primary key. Apply index for user_id to improve search.
    - It returns the list of hydrated tweets and the `next_cursor` for pagination.

### CQRS (Command Query Responsibility Segregation)

This "Fan-out on Write" model is a practical implementation of the **CQRS** pattern:

- **Command**: The write operation (`POST /tweet`) updates the write model (the `Tweets` table).
- **Query**: The read operation (`GET /timeline`) queries a separate, pre-calculated read model (the timeline lists in Redis), which is fully optimized for speed.

Source consulted: https://learn.microsoft.com/en-us/azure/architecture/patterns/cqrs

## 8. Scalability Strategy

The database is the primary bottleneck. We will apply the following strategies:

### Read Scaling

- **Read Replicas**: The primary PostgreSQL database will handle all write operations. Multiple read replicas will be created to handle read queries, such as hydrating tweets. The `GET /timeline` flow will primarily hit Redis, but the subsequent hydration query will be directed to the read replicas.
- **Eventual Consistency**: This architecture results in eventual consistency. There might be a slight delay before a new tweet appears on a follower's timeline.

Source consulted: https://www.educative.io/courses/grokking-the-system-design-interview/data-replication

### Write Scaling

- **Sharding**: To handle high write traffic and prevent the primary database from being a single point of failure, the `Tweets` table can be horizontally partitioned (sharded) based on `user_id`. A consistent hashing algorithm on the `user_id` can ensure an even distribution of data across multiple database shards.

Source consulted: https://www.educative.io/courses/grokking-the-system-design-interview/data-partitioning

### Service Separation

- The system can be split into separate microservices (e.g., `Users Service` and `Timeline Service`) to scale reads and writes independently. An API Gateway or load balancer would route `GET` requests to the Timeline service and `POST`/`PUT` requests to the Users service.

## 9. Stack

Go, Postgress, Redis, gomock (by Uber), Port and Adapters architecture.

Microservices architecture. For Troubleshooting recomend use Grafana or Datadog.

## 10. Next Iterations & Discussion Points

- **Service Splitting**: Formally separate the codebase into two distinct services:
- **Users Service**: Manages creating tweets, following users, etc.
- **Timeline Service**: Manages timeline generation and retrieval.
- **The "Celebrity" Problem**: For users with millions of followers, the "fan-out on write" can be expensive. A hybrid approach could be adopted where timelines for followers of celebrities are not pre-calculated and are instead merged at read time.
- **Worker Resilience**: Implement robust error handling, retries, and dead-letter queues for the timeline generation workers.

# Author

Lic. Renzo Mauro Ontivero

LinkedIn profile: https://www.linkedin.com/in/renzoontivero91/
