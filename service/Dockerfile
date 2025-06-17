FROM golang:1.24.1-alpine AS builder
WORKDIR /app

# install a C compiler, musl headers, and sqlite dev headers/libs
RUN apk add --no-cache \
      build-base   \
      sqlite-dev

# Copy go.mod/go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source
COPY . .

# Build with cgo enabled
ENV CGO_ENABLED=1 \
    GOOS=linux
RUN go build -o main .

# ---- final image ----

FROM alpine:latest
# for HTTPS calls
RUN apk add --no-cache ca-certificates

WORKDIR /root
# pull in the sqlite runtime if you did *not* static-link
# RUN apk add --no-cache sqlite-libs

COPY --from=builder /app/main .

RUN mkdir -p files
VOLUME ["/root/files"]
EXPOSE 3141

CMD ["./main", "start"]