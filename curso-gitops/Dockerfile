FROM golang:1.25.8-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o curso-gitops ./cmd/api

FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /home/appuser
COPY --from=builder --chown=appuser:appgroup /app/curso-gitops .
COPY --from=builder --chown=appuser:appgroup /app/frontend ./frontend
USER appuser
EXPOSE 8080
CMD ["./curso-gitops"]
