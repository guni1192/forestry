FROM golang:1.16.4-buster as builder

WORKDIR /src

RUN apt-get update -y && \
    apt-get install make && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /src

RUN make build

FROM gcr.io/distroless/base-debian10

COPY --from=builder /src/bin/forestry-server /forestry-server

ENTRYPOINT ["/forestry-server"]
