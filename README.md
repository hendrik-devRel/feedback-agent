# Collector Agent - Local Development Setup

## Prerequisites
- Docker and Docker Compose installed
- Go 1.21+ (for running the application)

- cahnges here 

## Quick Start

1. **Start the database:**
   ```bash
   cd collector-agent
   docker-compose up -d
   ```

2. **Wait for PostgreSQL to be ready:**
   ```bash
   # Check logs
   docker-compose logs -f postgres
   
   # Or wait for healthy status
   docker-compose ps
   ```

3. **Run migrations:**
   ```bash
   # Option A: Using psql from host (if installed)
   psql -h localhost -U postgres -d feedback -f migrations/202501150001_create_feedback.up.sql
   psql -h localhost -U postgres -d feedback -f migrations/202501150002_create_votes.up.sql
   
   # Option B: Using psql from Docker container
   docker-compose exec postgres psql -U postgres -d feedback -f /migrations/202501150001_create_feedback.up.sql
   docker-compose exec postgres psql -U postgres -d feedback -f /migrations/202501150002_create_votes.up.sql
   
   # Option C: Interactive shell
   docker-compose exec postgres psql -U postgres -d feedback
   # Then run: \i /migrations/202501150001_create_feedback.up.sql
   ```

4. **Verify setup:**
   ```bash
   docker-compose exec postgres psql -U postgres -d feedback -c "\dt"
   ```

## Database Connection Details

- **Host:** localhost
- **Port:** 5432
- **Database:** feedback
- **User:** postgres
- **Password:** password
- **Connection String:** `postgres://postgres:password@localhost:5432/feedback?sslmode=disable`

## Optional: pgAdmin GUI

Start pgAdmin for a web-based database management interface:
```bash
docker-compose --profile tools up -d
```

Access at: http://localhost:5050
- Email: admin@feedback.local
- Password: admin

Add server connection:
- Host: postgres (use container name, not localhost)
- Port: 5432
- Database: feedback
- User: postgres
- Password: password

## Common Commands

```bash
# Start services
docker-compose up -d

# Stop services (keeps data)
docker-compose down

# Stop and remove data
docker-compose down -v

# View logs
docker-compose logs -f postgres

# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d feedback

# Restart database
docker-compose restart postgres

# Check service status
docker-compose ps
```

## Testing the Schema

Run manual tests:
```bash
docker-compose exec postgres psql -U postgres -d feedback

# Inside psql:
-- Insert test feedback
INSERT INTO feedback (title, type, sentiment) 
VALUES ('Test Feature', 1, 1);

-- Insert test vote
INSERT INTO votes (feedback_id, user_id) 
VALUES (1, 100);

-- Query data
SELECT * FROM feedback;
SELECT * FROM votes;
```

## Troubleshooting

**Port 5432 already in use:**
```bash
# Check what's using the port
lsof -i :5432

# Either stop that service or change the port in docker-compose.yml:
# ports:
#   - "5433:5432"  # Use 5433 on host instead
```

**Database not ready:**
```bash
# Check health status
docker-compose ps

# Check logs for errors
docker-compose logs postgres
```

**Reset everything:**
```bash
docker-compose down -v
docker-compose up -d
# Re-run migrations
```
