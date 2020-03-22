FROM golang:1.14.1

ARG TOAD_VERSION=0.2.0

ENV GO111MODULE=on

RUN mkdir -p /app

RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Toad/releases/download/${TOAD_VERSION}/Toad_${TOAD_VERSION}_Linux_x86_64.tar.gz | tar xz

RUN rm LICENSE
RUN rm README.md
RUN mv Toad toad

RUN ./toad --get=release

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=2s --retries=5 --start-period=2s \
  CMD ./toad --get health

CMD ["./toad", "--port", "8080"]