FROM golang:1.17.3 as builder

# install cgo-related dependencies

WORKDIR  /usr/local/go/src/github.com/hamster-shared/hamster-gateway/

COPY . .

RUN set -eux; \
    go mod tidy ; \
    go build

FROM node:17 as builder-2

WORKDIR /usr/local/go/src/github.com/hamster-shared/hamster-gateway/frontend

COPY . .

RUN npm install ;\
    npm run build


FROM ubuntu:20.04

COPY --from=builder /usr/local/go/src/github.com/hamster-shared/hamster-gateway/hamster-gateway /usr/local/bin/

COPY --from=builder-2  /usr/local/go/src/github.com/hamster-shared/hamster-gateway/frontend/dist /usr/local/bin/frontend/

CMD hamster-gateway daemon