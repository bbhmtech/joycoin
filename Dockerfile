# build frontend
FROM node:lts-alpine as fe-builder
RUN npm install -g pnpm
WORKDIR /builder/frontend

COPY ./frontend/package.json ./frontend/pnpm-lock.yaml /builder/frontend
RUN pnpm install --frozen-lockfile

COPY ./frontend /builder/frontend
RUN pnpm build


# build go
FROM golang:alpine as go-builder
RUN apk add build-base
WORKDIR /builder/go

COPY go.mod go.sum /builder/go
RUN go mod download

COPY --from=fe-builder /builder/frontend/dist /builder/go/frontend/dist

COPY . /builder/go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
    go build -ldflags="-s -w" -gcflags="all=-N -l" -o ./bin/main ./cmd/main/main.go

# build runtime
FROM alpine:latest

WORKDIR /app/data
COPY --from=go-builder /builder/go/bin/main /app/bin/main

CMD ["/app/bin/main"]
