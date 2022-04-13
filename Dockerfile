FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /rabbitmq-mgn

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /rabbitmq-mgn /rabbitmq-mgn

USER nonroot:nonroot

ENTRYPOINT [ "/rabbitmq-mgn" ]
