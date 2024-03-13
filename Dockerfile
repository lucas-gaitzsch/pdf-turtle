FROM golang as build-service
WORKDIR /build
COPY . .
RUN go build -o pdf-turtle


FROM node:lts as build-playground
WORKDIR /app
COPY .pdf-turtle-playground/. .
RUN npm i
RUN npm run build


FROM chromedp/headless-shell:latest as runtime
WORKDIR /app
COPY --from=build-service /build/pdf-turtle /app/pdf-turtle
COPY --from=build-playground /app/dist /app/static-files/extern/playground

RUN apt-get -y update
RUN apt-get -y install media-types ca-certificates fonts-open-sans fonts-roboto fonts-noto-color-emoji
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

ENV LANG en-US.UTF-8
ENV LOG_LEVEL_DEBUG false
ENV LOG_JSON_OUTPUT false
ENV WORKER_INSTANCES 40
ENV PORT 8000
ENV SERVE_PLAYGROUND true

EXPOSE ${PORT}

RUN useradd -u 64198 app
USER app

ENTRYPOINT ["/app/pdf-turtle"]