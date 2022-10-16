
test: test-coverage benchmark-cpuprofile

test-coverage:
	mkdir -p coverage
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	xdg-open coverage/coverage.html

benchmark-cpuprofile:
	cd pkg/rxxr && go test -bench=. -cpuprofile cpu.prof
	go tool pprof -svg pkg/rxxr/cpu.prof > cpu.svg
	xdg-open cpu.svg