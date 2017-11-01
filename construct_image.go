package main

import (
	"time"
	"net/http"
	"mime"
	"path/filepath"
	"os"
	"io"
	"mime/multipart"
	"runtime"
	"github.com/disintegration/imaging"
	"image"
//	"log"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var widthThumbnail = 400
var widthPreview = 800


const imageIDLenght = 10

type Image struct {
	ID string
	UserID string
	Name string
	Location string
	Size int64
	CreatedAt time.Time
	Description string
}



var mimeExtensions = map[string]string{
	"image/png": ".png",
	"image/jpeg": ".jpg",
	"image/gif": ".gif",
}

func NewImage(user *User) *Image {
	return  &Image{
		ID: GenerateID("img", imageIDLenght),
		UserID: user.Id,
		CreatedAt: time.Now(),
	}
}

func (image *Image) StaticThumbnailRoute() string {
	return "/im/thumbnail/" + image.Location
}
func (image *Image) StaticPreviewRoute() string {
	return "/im/preview/" + image.Location
}
func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}
func(image *Image) ShowRoute() string {
	return "/image/" + image.ID
}

func (image *Image) CreateFromURL(imageURL string) error {
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer response.Body.Close()

	mimeType,_, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType

	}
	ext, valid := mimeExtensions[mimeType]
	if !valid {
		return errInvalidImageType
	}
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext
	savedFile, err := os.Create("./data/images/"+ image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}
	// Resize
	image.Size = size
	err = image.CreateResizedImages()
		if err != nil {
			return err
		}
	//
	return globalImageStore.Save(image)
}
func(image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	image.Name = headers.Filename
	types := mime.TypeByExtension(filepath.Ext(image.Name))
    //for _, v := range mimeExtensions {
    //
	//}
	_, valid := mimeExtensions[types]
	if !valid {
		return errInvalidImageType
	}
	image.Location = image.ID + filepath.Ext(image.Name)
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return  err
	}
	image.Size =  size
	// Resize
	err = image.CreateResizedImages()
	if err != nil {
		return err
	}
	//
     return globalImageStore.Save(image)
}
func (image *Image) DeleteImageFromHdd(file string) error {
	_ = os.Remove("data/images/" + image.Location)
	_ = os.Remove("data/images/preview/" + image.Location)
	_ = os.Remove("data/images/thumbnail/" + image.Location)
	return globalImageStore.Delete(file)
}

func (image *Image) CreateResizedImages() error {
	srcImage, err := imaging.Open("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	errorChan := make(chan error)
	go image.resizePreview(errorChan, srcImage)
	go image.resizeThumbnail(errorChan, srcImage)
	var e error
	for i := 0; i < 2; i++ {
		err := <-errorChan
		if err == nil {
			err = e
		}
	}
	return err
}
func (image *Image) resizeThumbnail(errorChan chan error, srcImage image.Image) {
dstImage := imaging.Thumbnail(srcImage, widthThumbnail, widthThumbnail, imaging.Lanczos)
destination := "./data/images/thumbnail/" + image.Location
errorChan <- imaging.Save(dstImage, destination)
}

func (image *Image) resizePreview(errorChan chan error, srcImage image.Image) {
	size := srcImage.Bounds().Size()
	ratio := float64(size.Y) / float64(size.X)
	targetHeight := int(float64(widthPreview) * ratio)
	dstImage := imaging.Resize(srcImage, widthPreview, targetHeight, imaging.Lanczos)
	destination := "./data/images/preview/" + image.Location
	errorChan <- imaging.Save(dstImage, destination)
}