# Feedback Agent - Architecture Overview

## Project Structure

```
feedback-agent/
├── app/
│   ├── main.go                    # Application entry point, wire everything together
│   │
│   ├── config/                    # Configuration management
│   │   └── config.go              # Load env vars, database config, etc.
│   │
│   ├── models/                    # Domain Models (What you have)
│   │   ├── entity/                # Core entities (Feedback, Vote)
│   │   ├── enum/                  # Enums (FeedbackType, Sentiment)
│   │   ├── request/               # API request DTOs
│   │   └── response/              # API response DTOs
│   │
│   ├── repository/                # Data Access Layer
│   │   ├── feedback_repository.go # Feedback CRUD operations
│   │   └── vote_repository.go     # Vote CRUD operations
│   │
│   ├── service/                   # Business Logic Layer
│   │   ├── feedback_service.go    # Feedback business logic
│   │   ├── sentiment_service.go   # Sentiment analysis
│   │   ├── classification_service.go # Actionability detection
│   │   └── vote_service.go        # Voting logic
│   │
│   ├── handler/                   # HTTP Handlers (API Layer)
│   │   ├── feedback_handler.go    # Feedback endpoints
│   │   └── health_handler.go      # Health check
│   │
│   ├── infrastructure/            # External Services & Infrastructure
│   │   ├── database/              # Database connection & setup
│   │   │   └── postgres.go        # PostgreSQL connection
│   │   ├── llm/                   # LLM integrations
│   │   │   ├── openai_client.go   # OpenAI integration
│   │   │   └── sentiment_analyzer.go # Sentiment analysis client
│   │   └── slack/                 # Slack integration (future)
│   │       └── slack_client.go
│   │
│   └── router/                    # HTTP Router Setup
│       └── router.go              # Gin routes, middleware
│
├── migrations/                    # Database migrations
├── docker-compose.yml
├── go.mod
└── README.md
```

## Architecture Layers

### 1. Domain Layer (`models/`)
**Purpose**: Core business entities and types
- **Entity**: `Feedback`, `Vote` - Core domain objects
- **Enum**: `FeedbackType`, `Sentiment` - Type-safe enums
- **Request/Response**: DTOs for API boundaries

**Responsibilities**:
- Define what data looks like
- No business logic
- No dependencies on other layers

### 2. Repository Layer (`repository/`)
**Purpose**: Abstract database operations
- `FeedbackRepository` interface: `Create()`, `GetByID()`, `GetAll()`, `Update()`
- `VoteRepository` interface: `Create()`, `GetByFeedbackID()`, `CountByFeedbackID()`

**Responsibilities**:
- Encapsulate all SQL queries
- Convert between database rows and entities
- Handle database-specific types (arrays, timestamps, etc.)

**Benefits**:
- Easy to swap database implementations
- Testable with mock repositories
- Single place for SQL queries

### 3. Service Layer (`service/`)
**Purpose**: Business logic and orchestration
- `FeedbackService`: Orchestrates feedback creation, validation, classification
- `SentimentService`: Analyzes text sentiment using LLM
- `ClassificationService`: Determines if feedback is actionable
- `VoteService`: Handles voting logic, prevents duplicates

**Responsibilities**:
- Business rules (validation, calculations)
- Orchestrate multiple repositories
- Call external services (LLM, Slack)
- Transform data between layers

**Benefits**:
- Business logic is centralized
- Easy to test (mock repositories/services)
- Can add new features without touching handlers

### 4. Handler Layer (`handler/`)
**Purpose**: HTTP request/response handling
- `FeedbackHandler`: HTTP endpoints for feedback operations
- `HealthHandler`: Health check endpoint

**Responsibilities**:
- Parse HTTP requests
- Call appropriate service methods
- Format HTTP responses
- Handle HTTP-specific errors

**Benefits**:
- Thin layer - delegates to services
- Easy to change HTTP framework (Gin → Echo, etc.)
- Can add multiple interfaces (REST, gRPC, GraphQL)

### 5. Infrastructure Layer (`infrastructure/`)
**Purpose**: External dependencies and setup
- `database/postgres.go`: Database connection setup
- `llm/openai_client.go`: OpenAI API client
- `config/config.go`: Environment variables, configuration

**Responsibilities**:
- Connect to external services
- Initialize dependencies
- Handle external API calls

### 6. Main Entry Point (`main.go`)
**Purpose**: Wire everything together
- Initialize database connection
- Create repositories
- Create services (with dependencies)
- Create handlers
- Setup router
- Start server

## Data Flow Example: Creating Feedback

```
HTTP Request (POST /api/feedback)
    ↓
Handler (feedback_handler.go)
    ├─ Parse JSON → CreateFeedbackRequest
    ├─ Validate request
    ↓
Service (feedback_service.go)
    ├─ Validate business rules
    ├─ Call SentimentService → Analyze sentiment
    ├─ Call ClassificationService → Check actionability
    ├─ Call FeedbackRepository → Save to DB
    ├─ Transform → Feedback entity
    ↓
Repository (feedback_repository.go)
    ├─ Build SQL query
    ├─ Execute INSERT
    ├─ Map DB row → Feedback entity
    ↓
Database (PostgreSQL)
    ↓
Return Feedback entity
    ↓
Handler formats response
    ↓
HTTP Response (JSON)
```

## Key Design Principles

### 1. Dependency Injection
- Services receive repositories via constructor
- Handlers receive services via constructor
- Makes testing easy (mock dependencies)

### 2. Interface-Based Design
```go
type FeedbackRepository interface {
    Create(ctx context.Context, feedback *entity.Feedback) error
    GetByID(ctx context.Context, id int) (*entity.Feedback, error)
    GetAll(ctx context.Context) ([]*entity.Feedback, error)
}

type FeedbackService interface {
    CreateFeedback(ctx context.Context, req *request.CreateFeedbackRequest) (*entity.Feedback, error)
}
```

### 3. Separation of Concerns
- Repository: Database operations only
- Service: Business logic only
- Handler: HTTP concerns only

### 4. Testability
- Each layer can be tested independently
- Mock interfaces for dependencies
- No global state

## Future Extensibility

### Adding Slack Integration
1. Add `infrastructure/slack/slack_client.go`
2. Add `service/slack_service.go` (uses SlackClient)
3. Add `handler/slack_handler.go` (webhook endpoint)
4. Wire into `main.go`

### Adding Voice Processing
1. Add `infrastructure/transcription/whisper_client.go`
2. Add `service/transcription_service.go`
3. Extend `FeedbackService` to handle audio → text → feedback

### Adding Analytics
1. Add `service/analytics_service.go`
2. Add `handler/analytics_handler.go`
3. Use existing repositories to query data

## Migration Path from Current Code

1. **Keep existing**: `models/entity/`, `models/enum/`
2. **Create**: `models/request/`, `models/response/`
3. **Extract**: Database code from `main.go` → `repository/`
4. **Move**: Business logic → `service/`
5. **Refactor**: HTTP handlers → `handler/`
6. **Setup**: `config/`, `infrastructure/`, `router/`

This architecture will scale as you add features without major refactoring!

