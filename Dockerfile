FROM golang:1.16.2

ARG TOAD_VERSION=0.2.4

ENV GO111MODULE=on
ENV IS_STATEFUL=false
ENV REDIS_HOST=localhost
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=

RUN mkdir -p /app

RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Toad/releases/download/${TOAD_VERSION}/Toad_${TOAD_VERSION}_Linux_x86_64.tar.gz | tar xz

RUN rm LICENSE
RUN rm README.md
RUN mv Toad toad

RUN ./toad --get release

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=2s --retries=5 --start-period=2s \
  CMD ./toad --get health

CMD ["./toad", "--port", "8080"]