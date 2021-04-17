# serialization

```bash
go get github.com/gogo/protobuf/protoc-gen-gofast
protoc --gofast_out=. my.proto # only encode
protoc --gofast_out=plugins=grpc:. my.proto # for grpc
```