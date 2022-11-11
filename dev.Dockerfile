FROM golang:1.19.3

RUN apt update && apt upgrade -y && \
	apt install -y git \
	make openssh-client

# all the code lives here. We gonna mount the volume here
WORKDIR /app

# copy the go.mod and go.sum files then download the dependencies
ADD go.mod go.sum /app/
RUN go mod download

# install the air tool
RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
	&& chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air