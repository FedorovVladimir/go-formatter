local-update-fmt:
	@echo "update fmt on your computer"
	@go build -o myfmt main.go
	@cp myfmt ~/go/bin/myfmt
	@rm -rf myfmt
	@echo "SUCCESS"
