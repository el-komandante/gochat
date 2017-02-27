FROM postgres:alpine
RUN apk update && \
apk add bash bash-completion
COPY init.sql /docker-entrypoint-initdb.d
