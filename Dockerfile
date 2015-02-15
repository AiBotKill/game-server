FROM ubuntu

RUN apt-get update -q

RUN DEBIAN_FRONTEND=noninteractive apt-get install -qy build-essential curl git

RUN curl -s https://storage.googleapis.com/golang/go1.3.src.tar.gz | tar -v -C /usr/local -xz

RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1

RUN mkdir -p /go/src/github.com/AiBotKill/game-server

ENV PATH /usr/local/go/bin:/go/bin:$PATH

ENV GOPATH /go

ADD . /go/src/github.com/AiBotKill/game-server

WORKDIR /go/src/github.com/AiBotKill/game-server

RUN go get

RUN go build .

CMD ./game-server