FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /starter ./cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stag

COPY --from=build-stage /starter /starter 

COPY .env .

USER nonroot:nonroot

CMD ["/starter"]