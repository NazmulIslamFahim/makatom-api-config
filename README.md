# Makatom API Config Service

A RESTful API service for managing configuration entities using MongoDB and Go.

## Features

- CRUD operations for configuration entities
- MongoDB integration with connection pooling
- Automatic input validation using the common handlers package
- Query parameter filtering and pagination
- Tenant-based data isolation
- Clean architecture with direct service-to-handler mapping

## API Endpoints

### Create Config
- **POST** `/config`
- **Body:**
```json
{
  "name": "database_config",
  "type": "database",
  "subtype": "postgresql",
  "tags": ["production", "database"],
  "tenant_id": "tenant123",
  "metadata": {
    "host": "localhost",
    "port": 5432,
    "database": "mydb"
  }
}
```

### Get All Configs
- **GET** `/configs?tenant_id=tenant123&type=database&limit=10&skip=0`
- **Query Parameters:**
  - `tenant_id` (required): Tenant identifier
  - `type` (optional): Filter by config type
  - `subtype` (optional): Filter by config subtype
  - `tag` (optional): Filter by tag
  - `limit` (optional): Number of results (default: 10)
  - `skip` (optional): Number of results to skip

### Get Config by ID
- **GET** `/config/get?id={id}`
- **Query Parameters:**
  - `id`: Config ObjectID

### Update Config
- **PUT** `/config/update?id={id}`
- **Query Parameters:**
  - `id`: Config ObjectID
- **Body:**
```json
{
  "name": "updated_database_config",
  "type": "database",
  "subtype": "mysql",
  "tags": ["staging", "database"],
  "metadata": {
    "host": "staging.example.com",
    "port": 3306,
    "database": "staging_db"
  }
}
```

### Delete Config
- **DELETE** `/config/delete?id={id}`
- **Query Parameters:**
  - `id`: Config ObjectID

## Configuration

Create a `.env` file in the root directory:

```env
API_PORT=:8080
ENVIRONMENT=development
DEBUG=true
MONGO_URI=mongodb://localhost:27017/makatom_config
MONGO_DATABASE=makatom_config
```

## Running the Service

1. Make sure MongoDB is running
2. Create the `.env` file with your configuration
3. Run the service:

```bash
./run.sh
```

## Data Model

### Config Entity
```go
type Config struct {
    ID             primitive.ObjectID     `bson:"_id" json:"id"`
    Name           string                 `bson:"name" json:"name"`
    Type           string                 `bson:"type" json:"type"`
    Subtype        string                 `bson:"subtype" json:"subtype,omitempty"`
    Tags           []string               `bson:"tags" json:"tags,omitempty"`
    TenantID       string                 `bson:"tenant_id" json:"tenant_id"`
    CreatedBy      string                 `bson:"created_by" json:"created_by"`
    LastUpdatedBy  string                 `bson:"last_updated_by" json:"last_updated_by"`
    Metadata       map[string]interface{} `bson:"metadata" json:"metadata,omitempty"`
    CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
    UpdatedAt      time.Time              `bson:"updated_at" json:"updated_at"`
}
```

## Architecture

- **Models**: Data structures and request/response types with validation tags
- **Services**: Business logic layer that returns `handlers.ServiceResponse`
- **Routes**: Direct connection between service functions and common handlers package
- **Common Package**: Leverages the existing handlers infrastructure for automatic validation and response handling

### Key Benefits

1. **No Redundant Handlers**: Service functions directly return `ServiceResponse`
2. **Automatic Validation**: Uses struct tags for input validation
3. **Clean Separation**: Business logic is separate from HTTP concerns
4. **Consistent Patterns**: Follows the same architecture as your other services

## Dependencies

- Go 1.23+
- MongoDB
- Common package (local dependency)
- go-playground/validator for input validation
- mongo-driver for MongoDB operations
