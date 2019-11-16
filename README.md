# frostexture

Provides a Go function to convert Frostpunk .texture files to .dds and .png. Reads all files in a directory with a DXT header, repairs them as .dds files, and optionally converts them to png.

## Caveats

Png conversion requires [ImageMagick](https://imagemagick.org/script/command-line-processing.php) utility.

## License

[MIT](LICENSE)