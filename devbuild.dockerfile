# Compile stage
FROM registry.gitlab.com/shitposting/golang:latest AS build-env

ADD . /go/src/gitlab.com/shitposting/autoposting-bot
WORKDIR /go/src/gitlab.com/shitposting/autoposting-bot

RUN go build
CMD ["./autoposting-bot", "-testing", "-polling"]


