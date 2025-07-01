# Driver Location API

Stores real-time latitude/longitude for drivers and returns the nearest
drivers to a rider’s position. It backs the **Matching API** in our
ride-hailing platform.

---

## Features

| Capability                         | Detail |
|-----------------------------------|--------|
| **Create drivers (bulk)**         | `POST /v1/drivers` (internal only) Bulk insert with adjustable batch size |
| **Geo-nearest search**            | `POST /v1/drivers/search` – GeoJSON point + radius + limit, returns drivers ordered by distance and includes `distanceMeters` |
| **Health-check**                  | `GET /v1/healthz` |
| **Authentication (internal)**     | `Authorization: Bearer <INTERNAL_API_KEY>` header |
| **2dsphere index**                | Automatic on `location` field |
| **Swagger / OpenAPI 3**           | `/swagger/index.html` (dev only) |
| **Tests**                         | Unit + Testcontainers integration |

---

##  Tech stack

* **Go 1.24**  
* **Echo v4** (HTTP framework)  
* **MongoDB 7.x** – Atlas or local  
* **swaggo/swag** & **echo-swagger** – docs  
* **testcontainers-go** – integration tests  

---

##  Quick start (dev)

```bash
# clone & bootstrap
git clone https://github.com/mrtuuro/driver-location-api.git
cd driver-location-api
touch .env                    # adjust system environments

# PORT=:<enter your port>
# MONGO_URI=mongodb+srv://<user>:<pass>@<host>/?retryWrites=true&w=majority
# DATABASE_NAME=<database-name>
# COLLECTION_NAME=<collection-name>
# SECRET_KEY=<enter your jwt secret key>
vim .env

# additional to see the Makefile commands
make help

# run the service
make run                      # regenerates Swagger, builds, starts on :8080

# seed initial driver CSV
make import                   # uses tools/importer/main.go
