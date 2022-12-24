FROM golang:1.19.3

RUN apt update && apt upgrade -y && \
	apt install -y git \
	make openssh-client

# all the code lives here. We gonna mount the volume here
WORKDIR /string-api

# copy the go.mod and go.sum files then download the dependencies
ADD go.mod go.sum /string-api/
RUN go mod download

# install the air tool
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
	&& chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

# install goose for db migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# will run from an entrypoint.sh file
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

CMD ["/entrypoint.sh"]