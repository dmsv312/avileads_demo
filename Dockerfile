FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
COPY --from=builder /app/myapp .
COPY --from=builder /app/http/ ./http/
COPY --from=builder /app/services/ ./services/
COPY --from=builder /app/static/ ./static/
COPY --from=builder /app/views/ ./views/


EXPOSE 8080
CMD ["./myapp"]