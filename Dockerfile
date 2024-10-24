###############
# BUILD IMAGE #
###############
FROM golang:alpine as builder
WORKDIR /tmp/build
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/tesla-http-api

###############
# FINAL IMAGE #
###############
FROM alpine:latest

ENV USER=tesla-http-api
ENV GROUPNAME=$USER
ENV UID=93761
ENV GID=93761
ENV APP_PATH=/var/lib/tesla-http-api

RUN apk --no-cache add tzdata ca-certificates
COPY --from=builder /tmp/build/app /
RUN mkdir -p $APP_PATH
RUN addgroup \
    --gid "$GID" \
    "$GROUPNAME" \
&&  adduser \
    --disabled-password \
    --gecos "" \
    --home "$APP_PATH" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$UID" \
    $USER
RUN chown -R $UID:$GID $APP_PATH
WORKDIR $APP_PATH
USER tesla-http-api
ENTRYPOINT ["/app"]
EXPOSE 8080