This repository is part of the study about the gRPC framework.

Install the last version of protobuf with you want to generate new proto buffer files.
https://github.com/protocolbuffers/protobuf/releases

Download the libraries
1 - go get -u google.golang.org/grpc
2 - go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate the files
1 -  protoc --proto_path=<path-to-protobuffer> --go_out=<path-to-protobuffer> --go_opt=paths=source_relative <path-to-protobuffer/.proto>
2 -  protoc --proto_path<path-to-protobuffer> --go-grpc_out=. <path-to-protobuffer/.proto>
