# Step 1: Modules caching
FROM golang:1.19-alpine3.16 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.19-alpine3.16 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN go build -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /app/data /data
COPY --from=builder /bin/app /app

CMD ["/app"]

EXPOSE 8080