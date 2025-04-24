## Project Overview

RSSAGG is a high-performance RSS feed aggregator built with concurrency in mind. This project aims to provide a reliable
and efficient way to collect and process RSS feeds from various sources.

### Dependencies

---

The project depends on the following packages:

```bash
github.com/go-chi/chi (v1.5.5)
github.com/google/uuid (v1.6.0)
github.com/joho/godotenv (v1.5.1)
github.com/lib/pq (v1.10.9)
```

---

### Core Logic

---
The core logic of the project is implemented in the utils/handlers.go file, which defines various API endpoints for
creating, reading, updating, and deleting (CRUD) operations on users, feeds, and follow feeds.

---

### API Endpoints

The project exposes several API endpoints, including:

| Endpoint         | Method | Description                                    | 
|------------------|--------|------------------------------------------------| 
| /users           | POST   | Create a new user account                      | 
| /feeds           | GET    | Retrieve all available feeds                   |
| /follow-feeds    | POST   | Create a new follow feed for a given user      | 
| /posts/{user_id} | GET    | Retrieve posts associated with a specific user |

---

### Installation and Setup

- RSSAGG requires GO version 1.23

---

#### Create directrory

```shell
mkdir project && cd project
```

---

#### Clone the repository - from console

```shell
git clone https://github.com/etrinque/rssagg.git
```

---

#### Navigate to project directory

```shell
cd /rssagg
```

---

#### Install dependencies

```shell
go mod tidy 
```

---

#### Configure environment variables in .env

- Create .env file in root of project

```shell
touch .env 
```

- Add these items to .env

```
PORT=["8080"]
 
CONNSTR="postgresql://[user[:password]@][host][:port][/dbname]"

// replace all within [ ] with your desired PORT and DB connection details;
```

---

#### Build

```shell
> cd main
> go build -o rss && ./rss
```

---

### Here are some examples of how to use the project:

- Create a new user account: curl -X POST -H "Content-Type: application/json" -d '{"name": "John
  Doe"}' http://localhost:8080/users
- Retrieve all available feeds: curl -X GET http://localhost:8080/feeds
- Create a new follow feed for a given user: curl -X POST -H "Content-Type: application/json" -d '{"feed_id": "
  1234567890"}' http://localhost:8080/follow-feeds

Licensing
This project is licensed under the MIT License. See the LICENSE file for details.

### API Endpoints & Function

| Route                    | Method | Description                                                            | 
|--------------------------|--------|------------------------------------------------------------------------|
| /v1/readiness            | GET    | Check server readiness                                                 |
| /v1/err                  | GET    | Error endpoint (for debugging)                                         |
| /v1/users                | POST   | Create a new user account                                              |
| /v1/users                | GET    | Retrieve list of users (requires auth)                                 |
| /v1/feeds                | GET    | Retrieve list of all feeds                                             |
| /v1/feeds                | POST   | Create a new feed (requires auth)                                      |
| /v1/posts                | GET    | Retrieve list of posts by user (requires auth)                         |
| /v1/feedfollow           | POST   | Create a new follow relationship between user and feed (requires auth) |
| /v1/feedfollow           | GET    | Retrieve list of follow relationships for user (requires auth)         |
| /v1/feedfollow/{feed_id} | DELETE | Delete follow relationship between user and feed (requires auth)       |

#### Note:

- authMiddleware is applied to routes that require authentication
- {feed_id} is a path parameter for the delete route