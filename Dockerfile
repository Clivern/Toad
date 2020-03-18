FROM golang:1.14.0

ARG TOAD_VERSION=0.0.2

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

CMD ["./toad", "--port", "8080"]