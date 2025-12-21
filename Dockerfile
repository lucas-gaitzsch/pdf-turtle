FROM golang:bullseye AS build-service
WORKDIR /build
COPY . .
RUN go build -o pdf-turtle


FROM node:lts AS build-playground
WORKDIR /app
COPY .pdf-turtle-playground/. .
RUN npm install
RUN npm run build


FROM chromedp/headless-shell:143.0.7499.170 AS runtime
WORKDIR /app
COPY --from=build-service /build/pdf-turtle /app/pdf-turtle
COPY --from=build-playground /app/dist /app/static-files/extern/playground

RUN apt-get -y update && \
    apt-get -y upgrade && \
    apt-get -y install ca-certificates fonts-open-sans fonts-roboto fonts-noto-color-emoji && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV LANG=en-US.UTF-8
ENV LOG_LEVEL_DEBUG=false
ENV LOG_JSON_OUTPUT=false
ENV WORKER_INSTANCES=40
ENV PORT=8000
ENV SERVE_PLAYGROUND=true

EXPOSE ${PORT}

ENTRYPOINT ["/app/pdf-turtle"]