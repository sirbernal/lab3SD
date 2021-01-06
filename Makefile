grpc:
	export GO111MODULE=on  
	go get github.com/golang/protobuf/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
	export PATH="$PATH:$(go env GOPATH)/bin"
	go get -u github.com/golang/protobuf/protoc-gen-go

protos: grpc
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/ClientService.proto
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/DNSService.proto
	protoc --proto_path=. --go_out=plugins=grpc:proto proto/AdminService.proto

runc:
	cd cliente && \
	go run cliente.go

rund0:
	cd dns0 && \
	go run dns.go

rund1:
	cd dns1 && \
	go run dns.go

rund2:
	cd dns2 && \
	go run dns.go

runb:
	cd broker && \
	go run broker.go            
runa:
	cd admin && \
	go run admin.go      
