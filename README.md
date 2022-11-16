# Image Convert

⚠️ refactoring in progress

[![ImageDup](https://github.com/kmulvey/imageconvert/actions/workflows/release_build.yml/badge.svg)](https://github.com/kmulvey/imageconvert/actions/workflows/release_build.yml) [![codecov](https://codecov.io/gh/kmulvey/imageconvert/branch/main/graph/badge.svg?token=XpJ5kCJzsn)](https://codecov.io/gh/kmulvey/imageconvert) [![Go Report Card](https://goreportcard.com/badge/github.com/kmulvey/imageconvert)](https://goreportcard.com/report/github.com/kmulvey/imageconvert) [![Go Reference](https://pkg.go.dev/badge/github.com/kmulvey/imageconvert.svg)](https://pkg.go.dev/github.com/kmulvey/imageconvert)

ImageConvert converts pngs and webps to jpeg and optionally compresses them with jpegoptim. 

## Requirements
- [ImageMagick](https://imagemagick.org/)
- [jpegoptim](https://github.com/tjko/jpegoptim)

## Run
`imageconvert -compress -log-file processed.log -dir /path/to/images`

print help:

`imageconvert -h`
