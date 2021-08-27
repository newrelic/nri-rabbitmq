ARG GOLANG_VERSION=1.16

FROM golang:${GOLANG_VERSION} as builder
WORKDIR /code
COPY go.mod .
RUN go mod download

COPY . ./
RUN go build -o ./bin/nri-rabbitmq src/; strip ./bin/nri-rabbitmq


FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /code/bin/nri-rabbitmq /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-rabbitmq
USER 1000
