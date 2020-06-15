ARG version=master

FROM golang:1.14.3-alpine

LABEL author=renjie email=penitence.rj@gmail.com version=${version}

COPY ddns /usr/local/bin/ddns

CMD ["ddns"]