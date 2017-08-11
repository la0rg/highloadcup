FROM golang:1.8.1
WORKDIR /go/src/github.com/la0rg/highloadcup/ 
COPY . .
RUN GOOS=linux go build -a -o app .
EXPOSE 80
CMD ["./app"]
