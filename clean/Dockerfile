FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG CMD=api
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/${CMD}

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/app /bin/app
COPY config/ /config/

ENTRYPOINT ["/bin/app"]
