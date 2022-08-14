FROM golang:1.18 as build_sales-api

ENV CGO_ENABLED 0
ARG BUILD_REF

COPY . /service

WORKDIR /service
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.16
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S sales && \
    adduser -u 1000 -h /service -G sales -S sales

COPY --from=build_sales-api --chown=sales:sales /service/utlimate-go /service/ultimate-go
WORKDIR /service
USER sales
CMD ["./ultimate-go"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="sales-api" \
      org.opencontainers.image.authors="Wojciech Wirkijowski <wojciech@wirkijowski.pl>" \
      org.opencontainers.image.source="https://github.com/wirkijowski/utlimate-go" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="wirkijowski.pl"