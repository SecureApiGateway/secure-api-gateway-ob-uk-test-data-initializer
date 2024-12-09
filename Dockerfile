FROM alpine:3.21.0

RUN apk --no-cache add curl libc6-compat gcompat
RUN addgroup -S sob-group && adduser -S sob-user -G sob-group

RUN mkdir -p "/opt/sob/config"

WORKDIR "/opt/sob"
COPY config config
COPY initialize .

RUN chmod 500 /opt/sob/initialize
RUN chown -R sob-user:sob-group /opt

USER sob-user
