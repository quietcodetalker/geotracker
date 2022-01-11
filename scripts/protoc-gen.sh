protoc \
--proto_path=. \
--go_out=./pkg/ --go_opt=paths=source_relative \
--go-grpc_out=./pkg/ --go-grpc_opt=paths=source_relative \
./api/proto/v1/location/location.proto
