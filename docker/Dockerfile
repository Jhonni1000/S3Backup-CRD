FROM alpine:3.19
RUN apk add --no-cache \
    postgresql-client \
    aws-cli
RUN addgroup -S backupgroup && adduser -S backupuser -G backupgroup
WORKDIR /backup
RUN chown backupuser:backupgroup /backup && chmod 700 /backup
USER backupuser
RUN pg_dump --version && aws --version