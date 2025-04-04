@echo off

cd /d "%~dp0"

echo "Starting Service 1 on port 8001..."
cd submission_service
go run cmd/main.go
