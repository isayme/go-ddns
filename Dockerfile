FROM golang:1.11.5-alpine AS builder

RUN apk update && apk add git

ARG APP_PKG
WORKDIR /go/src/${APP_PKG}

ENV GO111MODULE=on

COPY go.* ./
RUN go mod download

COPY . .

ARG APP_NAME
ARG APP_VERSION
ARG APP_BUILDTIME
ARG APP_GITSHA1
RUN CGO_ENABLED=0 go build -ldflags "-X ${APP_PKG}/src/util.Name=${APP_NAME} -X ${APP_PKG}/src/util.Version=${APP_VERSION} -X ${APP_PKG}/src/util.BuildTime=${APP_BUILDTIME} -X ${APP_PKG}/src/util.GitSHA1=${APP_GITSHA1}" -o /app/ddns main.go

FROM alpine
WORKDIR /app

COPY --from=builder /app/ddns /app/ddns

CMD ["/app/ddns"]
