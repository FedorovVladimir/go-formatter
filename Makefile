local-update-fmt:
	@echo "update fmt on your computer"
	@go build -o myfmt cmd/formatter/main.go
	@cp myfmt ~/go/bin/myfmt
	@rm myfmt
	@echo "SUCCESS"
