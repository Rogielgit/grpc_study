.PHONY: cert

cert:
	cd cert; ./cert.sh cd ..

run-server-poc:
	go run poc/server/main.go
run-client-poc:
	go run poc/client/main.go

run-protoc-poc:
	protoc --proto_path=poc/poc_proto --go-grpc_out=.  poc/poc_proto/poc.proto
	protoc --proto_path=poc/poc_proto  --go_out=poc/poc_proto  --go_opt=paths=source_relative poc/poc_proto/poc.proto


