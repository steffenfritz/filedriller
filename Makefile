all:
		go build
test:
		go test -v
		rm pronom.sig
bench:
		go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
		rm pronom.sig
		go tool pprof -web profile.out
clean:
	  rm -f profile.out
	  rm -f memprofile.out
		rm -f filedriller.test
