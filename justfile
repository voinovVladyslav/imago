# shows this list
default:
    @just --list


run-transformer:
    go run ./cmd/transformer/main.go


run-storage:
    go run ./cmd/storage/main.go

run-gateway:
    go run ./cmd/gateway/main.go

run-logger:
    go run ./cmd/logger/main.go

connect-to-storage-db:
    sqlite3 storage.sqlite3

recreated-storage-db:
    rm -f storage.sqlite3
    sqlite3 storage.sqlite3 < migrations/storage/0001_initial.sql

upload-file-test:
    curl -X POST http://localhost:8000/transform -F "file=@example.jpg" -F "email=admin@admin.com"

