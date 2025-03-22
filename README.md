# CampusNOC

A student portal project built using PostgreSQL and Go. This portal facilitates the submission and review of internship offers, allowing students, FPCs (Faculty Placement Coordinators), and HODs (Heads of Departments) to interact seamlessly.

## Features

- **Student Submissions**: Students can submit internship details and upload required documents.
- **FPC Reviews**: Faculty Placement Coordinators can review and approve/reject submissions.
- **HOD Reviews**: Heads of Departments can provide final approval for submissions.
- **Admin Management**: Admins can manage users, including FPCs and HODs.
- **Authentication**: Secure login for all roles with password hashing.

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

3. Create a `.env` file in the root directory and add the following:
   ```
   DB_URL=postgres://<username>:<password>@localhost:5432/student_portal?sslmode=disable
   ```

4. Run database migrations:
   ```bash
   psql -h localhost -U postgres -f migrations/migration.sql
   ```

5. Install Go dependencies:
   ```bash
   go mod tidy
   ```

## Running the Project

To run the project, execute:
```bash
go run cmd/main.go
```

## API Documentation

### Base URL
- `http://localhost:8001` (for student-related endpoints)
- `http://localhost:8002` (for admin, FPC, and HOD-related endpoints)

### Endpoints

#### Student
- **Submit Internship Details**: `POST /submit`
- **View Submissions**: `GET /submissions`

#### FPC
- **Login**: `POST /fpc/login`
- **Review Submissions**: `POST /fpc/fpc_reviews`

#### HOD
- **Login**: `POST /hod/login`
- **Review Submissions**: `POST /hod/hod_reviews`

#### Admin
- **Create Admin**: `POST /admin`
- **Create FPC**: `POST /admin/fpc`
- **Create HOD**: `POST /admin/hod`

Refer to the Postman collection in `docs/MUJ localhost.postman_collection.json` for detailed request/response examples.