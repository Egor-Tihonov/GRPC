gen:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     proto/model.proto

clear:
	rm proto/*.go
run:
	go run main.go