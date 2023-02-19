REPOPATH = github.com/kmulvey/imageconvert
BUILDS = imageconvert fixfilenames trimlog

build: 
	for target in $(BUILDS); do \
		go build -v -ldflags="-s -w" -o ./cmd/$$target ./cmd/$$target; \
	done

clean:
	for target in $(BUILDS); do \
    	go clean \
    	rm ./cmd/$$target
	done
