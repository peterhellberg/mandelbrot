mandelbrot
==========

The [Mandelbrot set](https://en.wikipedia.org/wiki/Mandelbrot_set) in Go.

[![Build Status](https://travis-ci.org/peterhellberg/mandelbrot.svg?branch=master)](https://travis-ci.org/peterhellberg/mandelbrot)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/mandelbrot)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/mandelbrot#license-mit)

## Command line tool

### Installation

```bash
go get -u github.com/peterhellberg/mandelbrot/cmd/mandelbrot
```

### Usage

```bash
Usage of mandelbrot:
  -f="mandelbrot.png": Filename of the image
  -w=640: Width of the image
  -h=480: Height of the image
  -n=30: Number of iterations to run
  -i="000000": Inside color
  -o="ffffff": Outside color
  -show=false: Show the generated image
```

## Examples

![Pink mandelbrot set](http://assets.c7.se/skitch/mandelbrot-20150409-223252.png)

## License (MIT)

*Copyright (C) 2015 [Peter Hellberg](http://c7.se/)*

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the "Software"),
> to deal in the Software without restriction, including without limitation
> the rights to use, copy, modify, merge, publish, distribute, sublicense,
> and/or sell copies of the Software, and to permit persons to whom the
> Software is furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included
> in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
> OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
> IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
> DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
> TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
> OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
