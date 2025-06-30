# 🗃️ Distributed Key-Value Store in Go

A lightweight distributed key-value store implemented in Go with master-slave replication, consistent sharding, and client redirection when hash-based routing detects a mismatch. 

---

## 🚀 Features

- ✅ **Written in Go** — fast and simple concurrency model using goroutines.
- 🔑 **Consistent Hashing** — deterministic shard placement based on key hashes.
- 🧭 **Redirection Logic** — automatic to the correct shard if the request hits the wrong node.
- 🧱 **Master-Slave Replication** — every shard has a read replica.
- 🗂️ **Pluggable Storage** — stores data on-disk via embedded database per shard.


---

## ⚙️ Getting Started

### 🔧 Prerequisites

- Go 1.20+
- Curl/Postman for testing

---

### 🛠️ Build

```bash
make build
```

▶️ Run Servers

```bash
make run-multi
```

This starts multiple shard nodes (Delhi, Mumbai, Hyderabad, Chennai), each with a master and replica.

📥 Set a Key

```bash
curl "http://127.0.0.2:8080/set?key=mykey&value=myvalue"
```

If the key doesn't belong to this shard based on consistent hashing, the server will respond with a redirect to the correct shard.

📤 Get a Key

```bash
curl "http://127.0.0.2:8080/get?key=mykey"
```

🧪 Load Testing

```bash
make load-test
```

This sends 1000 randomized key-value pairs to multiple shards to simulate distributed load.

🧹 Cleanup

```bash
make remove-all
```

Removes all .db files for a clean restart.


📁 Folder Structure

```
.
├── cmd/                # CLI and startup logic
├── config/             # TOML config files for sharding
├── db/                 # (Optional) database files
├── replication/        # Master-slave logic
├── web/                # HTTP server handlers
├── main.go             # Entry point
├── Makefile            # Build and orchestration
├── sharding.toml       # Shard mapping
```

🧠 Concepts

- Consistent Hashing ensures uniform key distribution across shards.
- Redirection makes any node entry-point possible — the router will redirect if needed.
- Replication provides read availability via replicas (eventually consistent).

🎊 Example

![example](https://github.com/wiptrax/distri-kv/blob/main/example.png)

🛣️ Future Improvements

- Leader election for fault tolerance
- gRPC API for more efficient communication
- Persistent WAL for better durability
- Dynamic shard rebalancing

📝 License

MIT — feel free to use, modify, and build on top of this project.
