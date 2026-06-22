FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/coi .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates && \
	addgroup -S app && \
	adduser -S app -G app

COPY --from=builder /out/coi /app/coi

USER app

EXPOSE 8080

CMD ["/app/coi"]
