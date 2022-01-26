protoc \
--proto_path=./api/proto/v1/location \
--go_out=. --go_opt=paths=import \
--go-grpc_out=. --go-grpc_opt=paths=import \
./api/proto/v1/location/*.proto

protoc \
--proto_path=./api/proto/v1/history \
--go_out=. --go_opt=paths=import \
--go-grpc_out=. --go-grpc_opt=paths=import \
./api/proto/v1/history/*.proto
