FROM golang:1.18.4-bullseye

WORKDIR /gitnotes

# install project dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy source code
COPY . .

# install git
RUN apt-get update
RUN apt-get install -y git=1:2.30.2-1

# run git init in the repo
RUN git config --global init.defaultBranch main
RUN git init

# set integration tests env vars
ENV GN_TEST_INTEAGRATION=TRUE
ENV GN_CODE_PATH=/gitnotes
ENV GN_CURRENT_BRANCH=main

# run tests
ENTRYPOINT make test