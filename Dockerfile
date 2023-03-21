# Step 1: Modules caching
FROM golang:1.20.2-alpine3.17 as build

WORKDIR /app
COPY go.mod  ./
COPY go.sum  ./

RUN go mod download

COPY *.go ./

RUN go build -o /bin/app



# Step 2: Builder
FROM golang:1.20.2-alpine3.17 as builder
WORKDIR /

COPY --from=build /bin/app /bin/app

EXPOSE 8080

ENTRYPOINT ["/bin/app"]