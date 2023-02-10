# The base go-image
FROM golang:alpine
 
# Create a directory for the app
RUN mkdir /ethereum-jsonrpc-to-rest
 
# Copy all files from the current directory to the app directory
COPY . /ethereum-jsonrpc-to-rest
 
# Set working directory
WORKDIR /ethereum-jsonrpc-to-rest

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

EXPOSE 8080

CMD [ "./server" ]
