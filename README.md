# 🚀 Rate-Limited API Service (Golang)

## 📌 Overview

This project implements a **rate-limited API service** in Go.
It provides endpoints to accept user requests and track per-user statistics while enforcing **rate limiting (5 requests per minute per user)**.

The system is designed with **clean architecture**, **concurrency safety**, and **extensibility** in mind.

---

## 🛠️ Tech Stack

* Go (Golang)
* Gin (HTTP framework)
* Redis (optional, for distributed rate limiting & storage)
* In-memory store (default fallback)

---

## ⚙️ Features

### ✅ Core Features

* `POST /v1/request` → Accepts user requests
* `GET /v1/stats` → Returns per-user request count
* Rate limiting: **5 requests per user per minute**
* Thread-safe in-memory implementation
* Redis-based distributed rate limiter (bonus)

### ✅ Advanced Features

* Sliding window rate limiting (accurate)
* Atomic operations using Redis Lua script
* Interface-based design (easy to swap implementations)
* Environment-based configuration

---

## 📂 Project Structure

```
.
├── cmd/main.go
├── controllers/
├── routes/
├── services/
├── stores/
├── ratelimiter/
├── models/
```

---

## 🧠 Design Decisions

### 1. Clean Architecture

* **Handler → Service → Store**
* Separation of concerns for better maintainability and testing

---

### 2. Rate Limiting Strategy

Two approaches were considered:

#### 🔹 Fixed Window (Redis INCR)

* Simple and fast
* ❌ Allows burst at window boundaries

#### 🔹 Sliding Window (Implemented)

* Stores timestamps of requests
* Removes expired entries
* Ensures accurate rate limiting

👉 **Chosen approach:** Sliding Window using Redis Sorted Set + Lua

---

### 3. Why Lua in Redis?

* Multiple operations (remove, count, insert) must be atomic
* Lua ensures execution as a **single atomic transaction**
* Prevents race conditions under concurrent requests

---

### 4. Storage Abstraction

Used interfaces:

```go
type StatsStore interface
type RateLimiter interface
```

👉 This allows:

* Switching between in-memory and Redis
* No changes in business logic

---

## 🔁 Memory vs Redis

| Feature     | In-Memory | Redis    |
| ----------- | --------- | -------- |
| Speed       | Fast      | Fast     |
| Concurrency | Mutex     | Atomic   |
| Distributed | ❌         | ✅        |
| Persistence | ❌         | Optional |

---

## 🚀 Setup & Run

### 1. Clone Repo

```bash
git clone https://github.com/dineshvis2022/rate-limited-api-service.git
cd api-service
```

---

### 2. Install Dependencies

```bash
go mod tidy
```

---

### 3. Create `.env`

```
PORT=8080
USE_REDIS=false (true: if want to store in redis)
REDIS_ADDR=localhost:6379
```

---

### 4. Run Application

```bash
go run main.go
```

---

## 🧪 API Usage

### 🔹 POST /v1/request

```bash
curl -X POST http://localhost:8080/v1/request \
-H "Content-Type: application/json" \
-d '{"user_id":"user1","payload":"data"}'
```

---

### 🔹 GET /v1/stats

```bash
curl http://localhost:8080/v1/stats
```

Response:

```json
{
  "user1": 3
}
```

---

## ⚡ Rate Limiting Behavior

* First 5 requests → ✅ Allowed
* 6th request within 1 minute → ❌ `429 Too Many Requests`

---

## 🧪 Concurrency Testing

```bash
seq 1 10 | xargs -n1 -P10 curl -X POST http://localhost:8080/v1/request \
-H "Content-Type: application/json" \
-d '{"user_id":"user1"}'
```

---

## ⚠️ Limitations

* In-memory store is not persistent
* Single Redis instance may become bottleneck
* No authentication or user validation
* No request payload storage (not required)

---

## 🚀 Future Improvements

* Add Redis clustering
* Add request queue (Kafka / RabbitMQ)
* Add metrics & monitoring (Prometheus)
* Add authentication & rate limiting per API key
* Persist request payloads if needed

---

## 💡 Key Highlights

* Concurrency-safe design using mutex / Redis atomic ops
* Clean separation of concerns
* Production-ready rate limiting logic
* Easily extensible architecture

---

## 🧑‍💻 Author

**Dinesh Vishwakarma**

---

## 📬 Submission

* GitHub Repo: `https://github.com/dineshvis2022/rate-limited-api-service.git`
* Completed within 24 hours as per assignment requirements

---
