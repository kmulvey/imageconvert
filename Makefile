REPOPATH = github.com/kmulvey/imageconvert
BUILDS := imageconvert fixfilenames trimlog

build: 
	for target in $(BUILDS); do \
		go build -v -ldflags="-s -w" -o ./cmd/$$target ./cmd/$$target; \
	done
