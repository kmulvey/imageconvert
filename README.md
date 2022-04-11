# Image Convert

[![ImageDup](https://github.com/kmulvey/imageconvert/actions/workflows/release_build.yml/badge.svg)](https://github.com/kmulvey/imageconvert/actions/workflows/release_build.yml) [![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://vshymanskyy.github.io/StandWithUkraine)

ImageConvert converts pngs and webps to jpeg and optionally compresses them with jpegoptim. 

## Requirements
- [ImageMagick](https://imagemagick.org/)
- [jpegoptim](https://github.com/tjko/jpegoptim)

## Run
`imageconvert -compress -log-file processed.log -dir /path/to/images`

print help:

`imageconvert -h`
