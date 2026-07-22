# shows this list
default:
    @just --list


run-transformer:
    go run ./cmd/transformer/main.go

run-storage:
    go run ./cmd/storage/main.go
