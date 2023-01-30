certs:
	cd cmd/ransim/.onos/config/certs; \
	./gen.sh; \
	cd ../../../../

gen:
	protoc --proto_path=api \
		  	--go_out=:api \
		  	--go-grpc_out=:api \
			api/*.proto

run:
	go run cmd/ransim/ransim.go

clean:
	rm api/*.go

.PHONY: certs run gen clean