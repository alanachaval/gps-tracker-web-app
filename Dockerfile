FROM golang:latest
ADD . $GOPATH/src/github.com/alanachaval/gps-tracker-web-app
WORKDIR $GOPATH/src/github.com/alanachaval/gps-tracker-web-app
RUN go get
RUN go install
WORKDIR $GOPATH/bin
CMD ["gps-tracker-web-app"]
