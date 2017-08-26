FROM golang:1.9 as builder
WORKDIR /go/src/github.com/la0rg/highloadcup/ 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/la0rg/highloadcup/app .
#COPY data.zip /tmp/data/
#COPY options.txt /tmp/data/
ENV GOGC=40
EXPOSE 80
CMD ["./app"]