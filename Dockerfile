FROM golang:1.24

WORKDIR /root

COPY ./ /root/

RUN go mod download

RUN go build -o collegeadminstration ./main.go

EXPOSE 5050

CMD ["./collegeadminstration"]