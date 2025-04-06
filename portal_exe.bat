@echo off

cd /d "%~dp0"

echo "Starting Portal Service on port 8002..."
cd portal_service
go run cmd/main.go