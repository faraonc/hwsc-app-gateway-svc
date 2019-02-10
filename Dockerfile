FROM golang:1.11.5
WORKDIR $GOPATH/src/github.com/hwsc-org/
RUN git clone https://github.com/hwsc-org/hwsc-app-gateway-svc.git
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
WORKDIR $GOPATH/src/github.com/hwsc-org/hwsc-app-gateway-svc
RUN dep ensure -v
RUN go install
ENTRYPOINT ["/go/bin/hwsc-app-gateway-svc"]
EXPOSE 50055
