@echo off

cd /d "%~dp0"

echo "Starting Service 2 on port 8002..."
cd portal_service
go run cmd/main.go