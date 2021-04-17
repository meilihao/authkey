#!/usr/bin/env bash

# protoc --go_out=plugins=grpc:. *.proto
protoc -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2/protobuf  -I=$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.2 -I=. --gofast_out=plugins=grpc:. *.proto

sed -i -E '/(.+) "google\/protobuf"/d' *.pb.go # 该行就是`import "github.com/golang/protobuf/proto"`