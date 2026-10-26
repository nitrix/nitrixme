# nitrixme

My personal website.

## Usage

Visit https://nitrix.me.

## License

This is free and unencumbered software released into the public domain. See the [UNLICENSE](UNLICENSE) file for more details.

## The portrait

The picture on the page is a still (`static/nitrixme.png`) that turns into a draggable orbit around me once its frames are in. The orbit is a clip (`static/nitrixme-orbit.mp4`, made with aigen's 360 Scene); `go generate` (it runs `tools/portrait`, which needs ffmpeg) cuts it into `static/nitrixme-orbit.bin`, one file of JPEGs the page loads in a single request. `static/nitrixme-old.jpg` is the previous picture, kept.
