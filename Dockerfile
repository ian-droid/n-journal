FROM golang:alpine

RUN apk add --no-cache \
 gcc libc-dev\
 git ;\
 go get github.com/ian-droid/njd ;\
 apk del gcc libc-dev git

EXPOSE 8086

VOLUME /journal.db /server.crt /server.key /ca.crt

ENTRYPOINT ["njd", "-db=/journal.db", "-cert=/server.crt", "-key=/server.key", "-ca=/ca.crt", "-docDir=/go/src/github.com/ian-droid/njd/"]
