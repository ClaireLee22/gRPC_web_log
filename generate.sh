# generate code
protoc web_log/web_log_pb/web_log.proto --go_out=plugins=grpc:.

# path setup
export GOPATH=$HOME/go
PATH=$PATH:$GOPATH/bin