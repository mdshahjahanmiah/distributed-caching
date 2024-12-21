# Distributed Caching with Golang and Redis

It demonstrates how to implement a distributed caching system using Golang and Redis, focusing on performance, scalability, and availability.

## Overview

This project showcases:

- **Distributed Caching**: Using Redis servers distributed across regions to cache data, reducing latency and database load.
- **Consistent Hashing**: Efficient data distribution across cache nodes for balanced utilization.
- **Cache Eviction Policies**: TTL for data freshness and conceptual LRU for memory management.
- **Region Awareness**: Optimized cache node selection based on the user's geographical location.

---

## Getting Started

### Prerequisites

- **Go** (version 1.16 or later)
- **Redis** (version 6 or later)
- **Docker** (optional for running Redis instances in containers)

---

### Setup

1. **Clone the Repository**
   ```bash
   git clone https://github.com/mdshahjahanmiah/distributed-caching.git
   cd distributed-caching

2.  ### Setup Redis:
You need at least three Redis instances for a distributed setup. Here's how to run them using Docker:

```bash
docker run --name redis-europe -p 6379:6379 -d redis
docker run --name redis-asia -p 6380:6379 -d redis
docker run --name redis-northamerica -p 6381:6379 -d redis
```
Alternatively, install Redis on your local machine or use a cloud service.

3. **Run the Application:**
   ```bash
   go run cmd/main.go
   ```
# Usage

The application simulates fetching book details from a distributed cache. It uses consistent hashing to determine which Redis instance should serve the data, with an option to consider the user's region for further optimization.

## Features

- **Consistent Hashing**: Ensures balanced data distribution across Redis nodes.
- **TTL Caching**: Data expiration to keep the cache relevant.
- **Region-Aware Caching**: Consideration of the user's region for potential latency reduction.

## Future Improvements

- Implement true Redis cluster support for more dynamic scaling.
- Enhance LRU simulation with actual eviction logic.
- Add more detailed performance metrics and dashboards for monitoring.
