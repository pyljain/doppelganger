test:
	go test ./... -v 

integration:
	go test ./... -v -cover --tags=integration