FROM golang:1.21.2

WORKDIR /testing

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN go build -o ./bin/app/app.exe -v ./cmd/app/main.go

CMD [ "./bin/app/app.exe"]