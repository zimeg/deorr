DEORR=deorr

.PHONY: build
build:
	go build -o build/$(DEORR)

.PHONY: test
test:
	go test -v . --coverprofile=coverage.out

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

.PHONY: benchmark
benchmark:
	go test -bench=. -ldflags="-s -w"

.PHONY: clean
clean:
	rm -rf build
	rm -rf coverage.{html,out}
	rm -rf records/*.out
