CGO_ENABLED=1 go build -a -installsuffix cgo -o monitor ./cmd/monitor
CGO_ENABLED=1 go build -a -installsuffix cgo -o api ./cmd/api