[![Build Status](https://travis-ci.org/dsoprea/go-perceptualhash.svg?branch=master)](https://travis-ci.org/dsoprea/go-perceptualhash)

## Overview

This is both a library and a command-line tool that can be used to produce perceptual hashes in Go. Perceptual hashes allow you to produce a hash that uniquely describes very similar images (e.g. the same image in different sizes or palettes). This is often used for deduplication.

This library is written in pure Go and is a translation of the algorithm at [http://blockhash.io](http://blockhash.io).


## Features

- We default to 16-bit hashes (=> 16 ^ 2 => 256 byte output hex-digest), but a different size can be passed.
- You can pass multiple files to the command-line tool.


## CLI Usage

Install and build (make sure you have defined GOPATH):

```
$ go get github.com/dsoprea/go-perceptualhash/command/go-perceptualhash
```

Run:

```
$ "${GOPATH}/bin/go-perceptualhash" --digest -f "${GOPATH}/src/github.com/dsoprea/go-perceptualhash/test_assets/20170618_155330.png"
1ffc3fff00fe000021ff7e3f0f8007c03fff1f8d0f9806003ffc3ff80f0400f0
```

The "--digest" parameter is used to print the digests without the file-paths.


## Programmatic Usage

Example:

```go
package main

import (
    ...

    "github.com/dsoprea/go-perceptualhash"

    ...
)

func main() {
    f, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    image, _, err := image.Decode(f)
    if err != nil {
        panic(err)
    }

    bh := blockhash.NewBlockhash(image, 16)
    hexdigest := bh.Hexdigest()

    // ...
}
```


## Tests

```
$ go test github.com/dsoprea/go-perceptualhash
```


## Notes

- Hashes of JPEG images will/may vary between different language implementations and/or image libraries due to a lack of specificity in JPEG regarding color conversions from YCbCr->RGB. If you wish to compare/benchmark implementations then use PNG.

- In practice, color and grayscale images will have different hashes.
