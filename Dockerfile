ARG version=master
ARG goVersion=1.14.3-alpine
ARG gitCommit=""
ARG buildStamp=""

FROM golang:${goVersion} as builder

RUN mkdir /app

ADD . /app

RUN go build -ldflags "-X 'main.gitCommit=${gitCommit}' -X 'main.buildStamp=${buildStamp}'" -o ddns pkg/cmd.go

FROM alpine

LABEL author=renjie email=penitence.rj@gmail.com version=${version}

COPY --from=builder /app/ddns /usr/local/bin/

CMD ["ddns"]