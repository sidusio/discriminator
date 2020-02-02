FROM golang:alpine as builder

# Application Directory
RUN mkdir /app
WORKDIR /app

# First handle dependencies as those probably are more stable than rest of codebase
COPY ./go.mod /app/
COPY ./go.sum /app/
RUN go mod download

# Copy source and build app
COPY . /app
RUN go build ./cmd/discriminator

FROM alpine

# Define environment variables
ENV DISCRIMINATOR_TEMPLATES_PATH=/templates
ENV DISCRIMINATOR_TEMPLATES_EXTENSION=.tmpl

ENV DISCRIMINATOR_CONTAINERS_LABEL=io.sidus.discriminator
ENV DISCRIMINATOR_INCLUDE_STOPPED_CONTAINERS=false

ENV DISCRIMINATOR_RUN_INTERVAL=5m

ENV DISCRIMINATOR_LOG_LEVEL=info
ENV DISCRIMINATOR_LOG_FORMAT=text

# Copy over the app from the builder image
COPY --from=builder /app/discriminator /discriminator

ENTRYPOINT ["/discriminator"]