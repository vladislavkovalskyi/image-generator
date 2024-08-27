package post

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GenerateTextRequest struct {
	APIKey string
	Text   string
	X      int
	Y      int
	R      uint8
	G      uint8
	B      uint8
	Size   float64
}

func loadFont(size float64) (font.Face, error) {
	const defaultFont = "fonts/Manrope-Regular.ttf"
	fontBytes, err := os.ReadFile(defaultFont)
	if err != nil {
		return nil, err
	}

	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func addLabel(img *image.RGBA, label string, x, y int, col color.Color, face font.Face) {
	point := fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}

	d.DrawString(label)
}

func parseGenerateTextRequest(context *gin.Context) (GenerateTextRequest, error) {
	var request GenerateTextRequest
	var err error

	request.APIKey = context.PostForm("api_key")
	request.Text = context.PostForm("text")

	request.X, err = strconv.Atoi(context.PostForm("x"))
	if err != nil {
		return request, err
	}
	request.Y, err = strconv.Atoi(context.PostForm("y"))
	if err != nil {
		return request, err
	}

	r, err := strconv.Atoi(context.PostForm("r"))
	if err != nil {
		return request, err
	}
	g, err := strconv.Atoi(context.PostForm("g"))
	if err != nil {
		return request, err
	}
	b, err := strconv.Atoi(context.PostForm("b"))
	if err != nil {
		return request, err
	}

	request.R = uint8(r)
	request.G = uint8(g)
	request.B = uint8(b)

	request.Size, err = strconv.ParseFloat(context.PostForm("size"), 64)
	if err != nil {
		return request, err
	}

	return request, nil
}

func decodeImage(fileType string, file multipart.File) (image.Image, error) {
	switch fileType {
	case "jpg", "jpeg":
		return jpeg.Decode(file)
	case "png":
		return png.Decode(file)
	default:
		return nil, errors.New("unsupported image format")
	}
}

func GenerateText(context *gin.Context) {
	file, header, err := context.Request.FormFile("image_data")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get image file"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
	img, err := decodeImage(ext, file)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image", "details": err.Error()})
		return
	}

	request, err := parseGenerateTextRequest(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	face, err := loadFont(request.Size)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load font"})
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	clr := color.RGBA{request.R, request.G, request.B, 255}
	addLabel(rgba, request.Text, request.X, request.Y, clr, face)

	switch ext {
	case "jpg", "jpeg":
		context.Writer.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(context.Writer, rgba, nil)
	case "png":
		context.Writer.Header().Set("Content-Type", "image/png")
		png.Encode(context.Writer, rgba)
	}
}
