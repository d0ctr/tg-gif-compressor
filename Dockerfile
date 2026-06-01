FROM golang

RUN apt update && apt upgrade -y
RUN apt install -y ffmpeg wget

WORKDIR /opt/app

RUN wget https://github.com/FallenProjects/tdlib-build/releases/download/v1.8.64/TDLib-tdjson-linux-x86_64.tar.gz &&\
        tar -xvf TDLib-tdjson-linux-x86_64.tar.gz &&\
        rm TDLib-tdjson-linux-x86_64.tar.gz

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot

ENV API_ID=
ENV API_HASH=
ENV TELEGRAM_TOKEN=
ENV TDLIB_PATH=/opt/app/libtdjson.so.1.8.64

CMD ["./bot"]
