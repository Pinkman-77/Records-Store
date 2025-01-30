📀 Records REST API
🎵 Overview
The Records REST API is a Go-based web service that allows users to manage a music records collection, including artists and albums. The API provides CRUD operations for artists and records, with authentication handled via a separate gRPC service.

📂 Project Structure

records-restapi/
├── cmd/                         # Application entry point
│   ├── main.go                  # Main application setup
│   └── configs/                 # Configuration files
│
├── pkg/
│   ├── handler/                  # HTTP request handlers
│   ├── repository/               # Database layer (PostgreSQL)
│   ├── service/                  # Business logic layer
│   ├── app.go                    # Application setup
│   ├── app/grpc/                 # gRPC authentication client
│   ├── grpc/auth/server.go       # Authentication gRPC service
│
├── proto/                        # Protocol Buffers for gRPC
│   ├── auth.proto                # Authentication service definition
├── Dockerfile                    # Docker configuration
├── docker-compose.yml            # Docker Compose configuration
├── go.mod                         # Go module dependencies
├── go.sum                         # Go module checksum
└── README.md                      # Project documentation
⚙️ Installation
1️⃣ Clone the Repository
git clone https://github.com/Pinkman-77/records-restapi.git
cd records-restapi
2️⃣ Install Dependencies
go mod tidy
3️⃣ Set Up the Database
Make sure PostgreSQL is installed and running. Create a database:

CREATE DATABASE rap_records_shop;
4️⃣ Configure Environment Variables
Create a .env file or update config.yaml with database credentials:

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "rap_records_shop"
5️⃣ Run the Application

go run cmd/main.go
🚀 API Endpoints
🎤 Artists
Method	Endpoint	Description
POST	/artists	Create a new artist
GET	/artists	Get all artists
GET	/artists/{id}	Get an artist by ID
PUT	/artists/{id}	Update an artist
DELETE	/artists/{id}	Delete an artist
📀 Records
Method	Endpoint	Description
POST	/records	Create a new record
GET	/records	Get all records
GET	/records/{id}	Get a record by ID
PATCH	/records/{id}	Update specific fields of a record
DELETE	/records/{id}	Delete a record
🛠️ Running Tests
The project includes unit tests for repository methods and API handlers.

Run Unit Tests
go test ./pkg/repository -v

🐳 Running with Docker
The project includes a Dockerfile and docker-compose.yml for containerized deployment.

1️⃣ Build the Docker Image

docker build -t records-restapi .

2️⃣ Run the Container

docker run -p 8080:8080 --config-file

3️⃣ Run with Docker Compose

docker-compose up --build

👨‍💻 Creator
Vitaliy aka Pinkman-77
📧 Email: ukvitaly7@gmail.com
🐙 GitHub: Pinkman-77


