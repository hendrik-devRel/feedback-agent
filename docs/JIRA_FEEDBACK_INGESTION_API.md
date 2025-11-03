# Jira Ticket: Implement Feedback Ingestion API

## Ticket Details

**Type:** Story  
**Priority:** High  
**Epic:** Feedback Collection System  
**Sprint:** Current

---

## Summary

Implement a RESTful API endpoint to ingest feedback entries via HTTP POST requests. This is the foundation for the feedback collection pipeline that will later support sentiment analysis, Slack integration, and automated processing.

---

## Description

As a developer/integrator, I want to submit feedback entries via a REST API endpoint so that feedback can be programmatically collected and stored in the database.

This endpoint will serve as the entry point for the feedback ingestion pipeline. It should accept all feedback fields (except auto-generated ones like ID, Votes, CreatedAt, UpdatedAt) and store them in the PostgreSQL database.

---

## User Story

**As a** system integrator  
**I want** to POST feedback entries to `/api/feedback`  
**So that** feedback can be collected and stored for later analysis

---

## Acceptance Criteria

### ✅ AC1: API Endpoint Implementation
- [ ] `POST /api/feedback` endpoint exists and is accessible
- [ ] Endpoint accepts JSON payloads with `Content-Type: application/json`
- [ ] Server runs on port `8080` by default
- [ ] Health check endpoint `GET /health` returns `{"status": "ok"}`

### ✅ AC2: Request Validation
- [ ] Request validates required fields: `title` (string), `type` (integer)
- [ ] Returns `400 Bad Request` with error message if validation fails
- [ ] Optional fields work correctly: `description`, `tags`, `sentiment`, `sentimentScore`
- [ ] `type` field accepts integers: `0` (Bug), `1` (Feature), `2` (General)
- [ ] `sentiment` field accepts integers: `0` (Neutral), `1` (Positive), `2` (Negative)
- [ ] `sentiment` defaults to `0` (Neutral) if not provided

### ✅ AC3: Database Integration
- [ ] Feedback is successfully inserted into PostgreSQL `feedback` table
- [ ] Auto-generated fields work correctly:
  - `id` is auto-generated (SERIAL)
  - `votes` defaults to `0`
  - `created_at` is set to current timestamp
  - `updated_at` is set to current timestamp
- [ ] All provided fields are stored correctly:
  - `title` → `title` (VARCHAR)
  - `description` → `description` (TEXT)
  - `type` → `type` (INT)
  - `tags` → `tags` (TEXT[])
  - `sentiment` → `sentiment` (INT)
  - `sentimentScore` → `sentiment_score` (DECIMAL)

### ✅ AC4: Response Format
- [ ] Returns `201 Created` status code on success
- [ ] Response body contains complete `Feedback` entity with all fields:
  ```json
  {
    "id": 1,
    "title": "...",
    "description": "...",
    "type": 1,
    "tags": ["..."],
    "sentiment": 0,
    "sentimentScore": null,
    "votes": 0,
    "createdAt": "2025-01-15T...",
    "updatedAt": "2025-01-15T..."
  }
  ```
- [ ] Returns `500 Internal Server Error` if database operation fails

### ✅ AC5: Enum JSON Support
- [ ] `FeedbackType` enum can unmarshal from JSON numbers (`0`, `1`, `2`)
- [ ] `Sentiment` enum can unmarshal from JSON numbers (`0`, `1`, `2`)
- [ ] Both enums validate that numbers are within valid range
- [ ] Enums return appropriate error messages for invalid values

### ✅ AC6: Code Quality
- [ ] Code follows Go best practices
- [ ] No linter errors
- [ ] Database connection is properly managed (defer db.Close())
- [ ] Error handling is implemented for all database operations
- [ ] PostgreSQL array types are handled correctly (tags field)

---

## Technical Requirements

### API Endpoints

#### `POST /api/feedback`
**Request Body:**
```json
{
  "title": "string (required)",
  "description": "string (optional)",
  "type": 0 | 1 | 2 (required),
  "tags": ["string"] (optional),
  "sentiment": 0 | 1 | 2 (optional, defaults to 0),
  "sentimentScore": 0.0-1.0 (optional)
}
```

**Success Response:** `201 Created`
```json
{
  "id": 1,
  "title": "...",
  "description": "...",
  "type": 1,
  "tags": ["..."],
  "sentiment": 0,
  "sentimentScore": null,
  "votes": 0,
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

**Error Response:** `400 Bad Request`
```json
{
  "error": "validation error message"
}
```

**Error Response:** `500 Internal Server Error`
```json
{
  "error": "Failed to create feedback"
}
```

#### `GET /health`
**Success Response:** `200 OK`
```json
{
  "status": "ok"
}
```

### Database Schema
- Uses existing `feedback` table from migration `202501150001_create_feedback.up.sql`
- Table structure:
  - `id` SERIAL PRIMARY KEY
  - `title` VARCHAR(200) NOT NULL
  - `description` TEXT NULL
  - `type` INT NOT NULL
  - `tags` TEXT[] NULL
  - `sentiment` INT NOT NULL DEFAULT 0
  - `sentiment_score` DECIMAL(3,2) NULL
  - `votes` INT NOT NULL DEFAULT 0
  - `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
  - `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

### Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/lib/pq` - PostgreSQL driver
- Go 1.21+

---

## Testing Requirements

### Manual Testing
- [ ] Test with Postman/curl:
  - [ ] Valid request with all fields
  - [ ] Valid request with only required fields
  - [ ] Invalid request (missing required fields)
  - [ ] Invalid request (invalid enum values)
  - [ ] Health check endpoint

### Test Cases
1. **Happy Path**: POST valid feedback → Returns 201 with complete feedback object
2. **Minimal Request**: POST with only `title` and `type` → Returns 201
3. **Full Request**: POST with all fields including tags and sentimentScore → Returns 201
4. **Validation Error**: POST without `title` → Returns 400
5. **Validation Error**: POST without `type` → Returns 400
6. **Invalid Enum**: POST with `type: 99` → Returns 400
7. **Database Error**: Test with database down → Returns 500
8. **Health Check**: GET /health → Returns 200

### Verification Steps
1. Start server: `go run app/main.go`
2. Verify database connection in logs
3. Send POST request via Postman
4. Verify response contains all fields
5. Check database directly (pgAdmin) to confirm data was saved
6. Verify auto-generated fields (id, votes, timestamps) are correct

---

## Implementation Checklist

- [x] Create `app/models/request/create_feedback.go` with request struct
- [x] Add `UnmarshalJSON` methods to `FeedbackType` enum
- [x] Add `UnmarshalJSON` methods to `Sentiment` enum
- [x] Implement `createFeedbackHandler` function
- [x] Setup Gin router with POST endpoint
- [x] Add health check endpoint
- [x] Handle PostgreSQL array types (tags field)
- [x] Set default sentiment to Neutral if not provided
- [x] Return complete Feedback entity in response
- [x] Implement proper error handling

---

## Files Changed

### New Files
- `app/models/request/create_feedback.go` - Request DTO

### Modified Files
- `app/main.go` - HTTP server setup and handler implementation
- `app/models/enum/feedback_type.go` - Added JSON unmarshaling support
- `app/models/enum/sentiment.go` - Added JSON unmarshaling support

---

## Definition of Done

- [ ] All acceptance criteria met
- [ ] Code reviewed and approved
- [ ] No linter errors
- [ ] Manual testing completed with Postman
- [ ] Data verified in database (pgAdmin)
- [ ] Documentation updated (if needed)
- [ ] PR description includes this ticket reference

---

## CodeRabbit Review Checklist

When reviewing this PR, CodeRabbit should verify:

1. **API Implementation**
   - [ ] POST endpoint exists at `/api/feedback`
   - [ ] Health check endpoint exists at `/health`
   - [ ] Server starts on port 8080

2. **Request Handling**
   - [ ] Validates required fields (`title`, `type`)
   - [ ] Handles optional fields correctly
   - [ ] Returns appropriate HTTP status codes

3. **Database Operations**
   - [ ] Inserts data into `feedback` table
   - [ ] Handles PostgreSQL array types (tags)
   - [ ] Auto-generates `id`, `votes`, `created_at`, `updated_at`

4. **Enum Support**
   - [ ] `FeedbackType` unmarshals from JSON numbers
   - [ ] `Sentiment` unmarshals from JSON numbers
   - [ ] Both enums validate input ranges

5. **Error Handling**
   - [ ] Returns 400 for validation errors
   - [ ] Returns 500 for database errors
   - [ ] Error messages are descriptive

6. **Code Quality**
   - [ ] No linter errors
   - [ ] Proper error handling
   - [ ] Database connections managed correctly
   - [ ] Follows Go best practices

---

## Notes

- This is the foundation for the feedback ingestion pipeline
- Future enhancements will add:
  - Sentiment analysis integration
  - Slack bot integration
  - Voice transcription support
  - Background processing agents
- The current implementation is minimal and focused on core functionality

---

## Related Tickets

- Future: Add sentiment analysis service
- Future: Add Slack integration
- Future: Add GET endpoints for retrieving feedback
- Future: Add voting functionality

