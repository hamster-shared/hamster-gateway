FROM node:17 as builder-2

WORKDIR /usr/local/go/src/github.com/hamster-shared/hamster-gateway/frontend

COPY ./frontend .


RUN npm install ;\
    npm run build


FROM golang:1.17.3 as builder

WORKDIR  /usr/local/go/src/github.com/hamster-shared/hamster-gateway/

COPY . .

ENV GO111MODULE  on
ENV GOPROXY https://goproxy.cn

RUN set -eux; \
    go mod tidy ; \
    go build



FROM ubuntu:20.04

WORKDIR /home/app

COPY --from=builder /usr/local/go/src/github.com/hamster-shared/hamster-gateway/hamster-gateway /home/app/

COPY --from=builder-2  /usr/local/go/src/github.com/hamster-shared/hamster-gateway/frontend/dist /home/app/frontend/

CMD hamster-gateway daemon