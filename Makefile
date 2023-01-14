REPOPATH = github.com/kmulvey/imageconvert
BUILDS := imageconvert

build: 
	for target in $(BUILDS); do \
		go build -v -x -ldflags="-s -w" -o ./cmd/$$target ./cmd/$$target; \
	done
