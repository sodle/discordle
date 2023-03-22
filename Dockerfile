FROM golang:1.19

WORKDIR /app

COPY discordle discordle
COPY word_bank word_bank
COPY go.mod .
COPY go.sum .
COPY main.go .

RUN go get discordle
RUN go install

ENTRYPOINT [ "discordle" ]
