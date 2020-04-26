.PHONY: docs test
docs:
	@docker run -v $$PWD/:/docs pandoc/latex -f markdown /docs/README.md -o /docs/build/output/README.pdf

test:
	@echo "***** Hello, Unit Tests! *****\n"
	@go test -v
	@echo "\n***** Hello, Integration Tests *****\n"
	@go test -v -tags=integration ./test/integration
