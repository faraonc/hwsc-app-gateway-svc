FROM golang:1.12.6
WORKDIR $GOPATH/
RUN git clone https://github.com/hwsc-org/hwsc-app-gateway-svc.git
WORKDIR $GOPATH/hwsc-app-gateway-svc
RUN go mod download
RUN go install
ENTRYPOINT ["/go/bin/hwsc-app-gateway-svc"]
EXPOSE 50055
