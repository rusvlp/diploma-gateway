from golang:1.26-bookworm

workdir app

COPY . .

run go mod download

cmd ["go", "run", "main.go"]