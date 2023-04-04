FROM golang:alpine 

RUN mkdir /app
WORKDIR /app
COPY go.mod .
RUN go mod tidy
COPY . .
RUN go build main.go

EXPOSE 8000
CMD [ "./main" ]