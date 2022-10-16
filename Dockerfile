FROM golang:latest

RUN mkdir /app 


ARG configurations configurations
ARG port  
COPY . /app  
# Replacing the configurations folder files with needed configurations 
COPY ${configurations} /app/configurations
WORKDIR /app 

RUN export GO111MODULE=on  
RUN go mod tidy 
EXPOSE ${port}
CMD go run ./main
