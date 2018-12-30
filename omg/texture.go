package omg

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/pkg/errors"
)

type texture struct {
	path          string
	textureHandle uint32
}

func newTexture(path string) *texture {
	return &texture{
		path: path,
	}
}

// load reads the texture from disk and sends it to opengl
func (t *texture) load() error {
	//	f, err := ioutil.
	f, err := os.Open(t.path)
	if err != nil {
		return errors.Wrap(err, "failed to open file with texture")
	}
	img, err := png.Decode(f)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to decode png texture %v", t.path))
	}

	img = imaging.FlipV(img) // flip the image because opengl likes images to be upside down

	rgb := image.NewRGBA(img.Bounds()) // we need to convert the image to rgba from whichever format it was loaded
	draw.Draw(rgb, rgb.Bounds(), img, image.Pt(0, 0), draw.Src)
	gl.GenTextures(1, &t.textureHandle)
	gl.BindTexture(gl.TEXTURE_2D, t.textureHandle)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // image downsizing filter
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // upscaling filter
	size := rgb.Rect.Size()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.SRGB_ALPHA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgb.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return nil
}

func (t *texture) bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.textureHandle)
}

func (t *texture) dispose() {
	gl.DeleteTextures(1, &t.textureHandle)
}
