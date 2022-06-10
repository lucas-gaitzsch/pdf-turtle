FROM golang:bullseye as build
WORKDIR /build
COPY . .
RUN go build -o pdf-turtle


FROM chromedp/headless-shell:latest as runtime
WORKDIR /app
COPY --from=build /build/pdf-turtle /app/pdf-turtle

ENV LOG_LEVEL_DEBUG false
ENV LOG_JSON_OUTPUT false
ENV WORKER_INSTANCES 40
ENV PORT 8000

EXPOSE ${PORT}

ENTRYPOINT ["/app/pdf-turtle"]