all:
		go build
test:
		go test -v
bench:
		go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
		go tool pprof -web profile.out
