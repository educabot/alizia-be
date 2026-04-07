# Build stage
FROM golang:1.26.1-alpine AS builder
RUN apk --no-cache add git
ARG GITHUB_TOKEN
RUN echo "machine github.com login x-token password ${GITHUB_TOKEN}" > /root/.netrcENV GOPRIVATE=github.com/educabot/*
ENV GOPRIVATE=github.com/educabot/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /alizia-api ./cmd

# Run stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /alizia-api .
COPY db/migrations ./db/migrations
EXPOSE 8080
CMD ["./alizia-api"]
