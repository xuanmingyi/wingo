package wallpaper

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	"strconv"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func SetWallpaper(X *xgbutil.XUtil, image *xgraphics.Image) error {
	image.XSurfaceSet(X.RootWin())
	image.XDraw()
	image.XPaint(X.RootWin())
	return nil
}

func NewColorImage(X *xgbutil.XUtil, color string) *xgraphics.Image {
	ximg := xgraphics.New(X, image.Rect(0, 0, 1280, 768))

	r, _ := strconv.ParseUint(color[1:3], 16, 8)
	g, _ := strconv.ParseUint(color[3:5], 16, 8)
	b, _ := strconv.ParseUint(color[5:7], 16, 8)

	ximg.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BGRA{R: uint8(r), G: uint8(g), B: uint8(b)}
	})
	return ximg
}

func SetColorWallpaper(X *xgbutil.XUtil, color string) {
	image := NewColorImage(X, color)
	if err := SetWallpaper(X, image); err != nil {
		panic(err)
	}
}

func NewImage(filename string) (img image.Image, err error) {
	data, err := ioutil.ReadFile(filename[7:])
	if err != nil {
		return nil, err
	}
	img, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SetFileWallpaper(X *xgbutil.XUtil, filename string) {
	img, err := NewImage(filename)
	if err != nil {
		panic(err)
	}
	ximg := xgraphics.New(X, image.Rect(0, 0, 1920, 1080))

	ximg.For(func(x, y int) xgraphics.BGRA {
		r, g, b, a := color.RGBAModel.Convert(img.At(x, y)).RGBA()

		return xgraphics.BGRA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	})

	if err = SetWallpaper(X, ximg); err != nil {
		panic(err)
	}
}
