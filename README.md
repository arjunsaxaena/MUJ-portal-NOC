# MUJ_AMG

A student portal project built using PostgreSQL and Go.

## Prerequisites

Ensure you have the following installed:
- PostgreSQL 17.2
- Go 1.23.2

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd MUJ_AMG
   ```

2. Set up the PostgreSQL database:
   ```bash
   psql -h localhost -U postgres -d student_portal
   ```

3. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Running the Project

To run the project, execute:
```bash
go run cmd/main.go
```