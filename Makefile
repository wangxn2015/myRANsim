certs:
	cd cmd/ransim/.onos/config/certs; \
	./gen.sh; \
	cd ../../../../

run:
	go run cmd/ransim/ransim.go
.PHONY: certs run