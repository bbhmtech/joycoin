# build frontend
FROM node:16-alpine as fe-builder
RUN npm install -g pnpm
WORKDIR /builder/frontend

COPY package.json pnpm-lock.yaml /builder/frontend
RUN pnpm install --frozen-lockfile

COPY ./frontend /builder/frontend
RUN pnpm build


# build go
FROM golang:1.15.2-alpine as go-builder
RUN apk add build-base
WORKDIR /builder/go

COPY go.mod go.sum /builder/go
RUN go mod download

COPY --from=fe-builder /builder/frontend/dist /builder/go/frontend/dist

COPY . /builder/go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w" -gcflags="all=-N -l" -o ./bin/main ./cmd/main/main.go


# build runtime
FROM alpine

WORKDIR /app
COPY --from=go-builder /builder/go/bin/main /app/main

CMD ["/app/main"]
