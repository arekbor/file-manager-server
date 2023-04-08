FROM golang:1.20

LABEL maintainer="file-manager-server-user"
RUN groupadd -r -g 1600 file-manager-server-user
RUN useradd -r -g 1600 -u 1500 file-manager-server-user

RUN chsh -s /usr/sbin/nologin root

WORKDIR /app

#COPY --chown=1500:1600 ./go.mod ./
#COPY --chown=1500:1600 ./go.sum ./

COPY --chown=1500:1600 . ./

RUN chown -R 1500:1600 /app

RUN go mod download

RUN go build -o /file-manager-server ./cmd/main.go

RUN cd /

RUN chmod -R 700 /media
RUN chown -R 1500:1600 /media

USER file-manager-server-user

CMD [ "/file-manager-server" ]
