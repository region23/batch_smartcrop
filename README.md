## Batch images resizer with smart crop (jpg/png) 

Based on [smartcrop.go](https://github.com/muesli/smartcrop) library by Christian Muehlhaeuser.

### Installation

Make sure you have a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

To install batch_smartcrop, simply run:

    git clone git://github.com/pavlik/batch_smartcrop.git
    cd batch_smartcrop && go build && go install

You can specify next params:
- **path** to images folder. Or you can run *batch_smartcrop* from images directory
- **prefix** for new thumbnail images
- **width** of thumbnail image
- and/or **height** of thumbnail image


Type for help
```
batch_smartcrop --help
```