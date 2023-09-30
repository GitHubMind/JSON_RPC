FROM golang:1.18-bullseye

WORKDIR /app

COPY . ./
ENV GOPROXY="https://goproxy.io"

RUN go build main.go

LABEL Description="This simple test by rpc serices"
CMD /app/main
