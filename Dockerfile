# syntax=docker/dockerfile:1

### Build the server
FROM golang:1.19-alpine AS serverbuild

WORKDIR /app
COPY . ./

# Compile binary
WORKDIR /app/src
RUN go mod download && go mod verify
RUN go build -o /server .

### Final image ###
FROM alpine:latest
WORKDIR /app
COPY --from=serverbuild /server ./
COPY --from=serverbuild /app/src/static ./static
ENV GOMEMLIMIT=100MiB
ENTRYPOINT ["./server"]
