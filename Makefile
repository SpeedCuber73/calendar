gen.proto:
	protoc --go_out=plugins=grpc:. api/event.proto