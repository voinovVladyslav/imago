# shows this list
default:
    @just --list


run-transformer:
    go run ./cmd/transformer/main.go


run-storage:
    go run ./cmd/storage/main.go


connect-to-storage-db:
    sqlite3 storage.sqlite3


recreated-storage-db:
    rm -f storage.sqlite3
    sqlite3 storage.sqlite3 < migrations/storage/0001_initial.sql
