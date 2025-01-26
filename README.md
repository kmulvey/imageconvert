# Image Convert

[![Build](https://github.com/kmulvey/imageconvert/actions/workflows/build.yml/badge.svg)](https://github.com/kmulvey/imageconvert/actions/workflows/build.yml) [![Release](https://github.com/kmulvey/imageconvert/actions/workflows/release.yml/badge.svg)](https://github.com/kmulvey/imageconvert/actions/workflows/release.yml) [![codecov](https://codecov.io/gh/kmulvey/imageconvert/branch/main/graph/badge.svg?token=XpJ5kCJzsn)](https://codecov.io/gh/kmulvey/imageconvert) [![Go Report Card](https://goreportcard.com/badge/github.com/kmulvey/imageconvert/v2)](https://goreportcard.com/report/github.com/kmulvey/imageconvert/v2) [![Go Reference](https://pkg.go.dev/badge/github.com/kmulvey/imageconvert/v2.svg)](https://pkg.go.dev/github.com/kmulvey/imageconvert/v2/pkg/imageconvert)

ImageConvert converts pngs and webps to jpeg and optionally compresses them with jpegoptim. 

## Requirements
- [ImageMagick](https://imagemagick.org/)
- [jpegoptim](https://github.com/tjko/jpegoptim)

## Run
```
imageconvert -compress -log-file processed.log -dir /path/to/images
```
```
imageconvert -compress -threads 1 -depth 2 -resize-threshold "2560x1440" -resize-size "5120x2880" -processed-file processed.log -path /path/to/images -watch
```
```
trimlog -log-file processed.log
```

print help:

`imageconvert -h`
