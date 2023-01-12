ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine AS builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod  download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
        -installsuffix 'static' \
        -o /platzi-rest-go

FROM scratch AS runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

# COPY .env ./

COPY --from=builder /platzi-rest-go /platzi-rest-go

EXPOSE 5050

ENTRYPOINT ["/platzi-rest-go"]

# docker build . -t go-rest-ws
# docker run  -p 5050:5050 go-rest-ws