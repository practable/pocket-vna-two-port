#/bin/bash
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
unzip protoc-3.15.8-linux-x86_64.zip -d $HOME/.local
rm protoc-3.15.8-linux-x86_64.zip 
export PATH="$PATH:$HOME/.local/bin"
echo "installed to "$PATH:$HOME/.local/bin" - consider updating your path to reflect this
pip3 install grpcio-tools
#go get github.com/golang/protobuf/protoc-gen-go #was this needed?
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
