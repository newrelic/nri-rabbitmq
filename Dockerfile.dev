ARG GOLANG_VERSION=1.18

FROM golang:${GOLANG_VERSION} as builder

WORKDIR /code
COPY . ./
RUN make compile; strip ./bin/nri-rabbitmq


FROM newrelic/infrastructure-bundle:latest
COPY --from=builder /code/bin/nri-rabbitmq /var/db/newrelic-infra/newrelic-integrations/bin/nri-rabbitmq
