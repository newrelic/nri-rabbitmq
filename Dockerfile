FROM golang:1.9 as builder
COPY . /go/src/github.com/newrelic/nri-rabbitmq/
RUN cd /go/src/github.com/newrelic/nri-rabbitmq && \
    make && \
    strip ./bin/nri-rabbitmq

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-rabbitmq/bin/nri-rabbitmq /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-rabbitmq
COPY --from=builder /go/src/github.com/newrelic/nri-rabbitmq/rabbitmq-definition.yml /nri-sidecar/newrelic-infra/newrelic-integrations/definition.yml
USER 1000
