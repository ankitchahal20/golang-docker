FROM golang:latest

## create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app
## copy everything in the root directory of app
ADD . /app
## change DIR
WORKDIR /app
## run go build to compile the binary 
## -o flag should be same in CMD which is to run the binary "main"
RUN go build -o main .
EXPOSE 8080
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]



