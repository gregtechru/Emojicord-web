FROM golang:stretch as builder

RUN go get -u github.com/labstack/echo \
              github.com/dgrijalva/jwt-go \
              github.com/go-sql-driver/mysql
RUN go get golang.org/x/oauth2 \
           github.com/jmoiron/sqlx \
           github.com/bwmarrin/discordgo

WORKDIR /go/src/github.com/Team-Fruit/Emojicord-web/web/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine

ENV DOCKERIZE_VERSION v0.6.1

RUN apk --no-cache add ca-certificates

WORKDIR /root/

RUN apk --no-cache add --virtual=wget wget \
    && wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
	&& apk del wget

COPY --from=builder /go/src/github.com/Team-Fruit/Emojicord-web/web/app .

ENTRYPOINT ["dockerize", "-wait", "tcp://db:3306"]
CMD ["./app"]
