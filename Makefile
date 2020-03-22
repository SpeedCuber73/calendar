gen.proto:
	protoc --go_out=plugins=grpc:pkg/calendar api/api.proto