package controllers

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/rwcarlsen/goexif/exif"
)

type MagicImage struct {
	MultipartForm  *multipart.Form
	AllowExt       []string
	MaxFileSize    int
	Required       bool
	DuplicateImage []uint64
	MinFileInSlice int
	MaxFileInSlice int
	Files          []*multipart.FileHeader
	FileNames      []string
}

func New(MultipartForm *multipart.Form) *MagicImage {
	return &MagicImage{
		MultipartForm: MultipartForm,
		AllowExt:      []string{"jpeg", "png"},
		MaxFileSize:   4 << 20, // 4 MB
		Required:      true,
	}
}

func (magic *MagicImage) inAllowExt(key string) (string, bool) {
	for _, value := range magic.AllowExt {
		if key == "image/"+value {
			return value, true
		}
	}
	return "", false
}

func (magic *MagicImage) isValidImage(fhs *multipart.FileHeader) (string, bool) {
	file, err := fhs.Open()
	if err != nil {
		return "", false
	}
	defer file.Close()

	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err = file.Read(buff); err != nil {
		return "", false
	}

	return magic.inAllowExt(http.DetectContentType(buff))
}

func (magic *MagicImage) appendImageHash(ext string, fhs *multipart.FileHeader) {
	file, _ := fhs.Open()
	defer file.Close()

	switch ext {
	case "jpeg":
		img, _ := jpeg.Decode(file)
		hash, _ := goimagehash.DifferenceHash(img)
		magic.DuplicateImage = append(magic.DuplicateImage, hash.GetHash())
	case "png":
		img, _ := png.Decode(file)
		hash, _ := goimagehash.DifferenceHash(img)
		magic.DuplicateImage = append(magic.DuplicateImage, hash.GetHash())
	default:
		panic("Cannot find extension.")
	}
}

func (magic *MagicImage) findDuplicateImage() bool {
	visited := map[uint64]bool{}

	for _, item := range magic.DuplicateImage {
		if visited[item] == true {
			return true
		} else {
			visited[item] = true
		}
	}

	return false
}

// validator for single image.
func (magic *MagicImage) ValidateSingleImage(key string) error {
	fhs := magic.MultipartForm.File[key]

	if !magic.Required && len(fhs) == 0 {
		return nil
	}

	if magic.Required && len(fhs) == 0 {
		return errors.New("Image is required.")
	}

	// validate image
	if _, valid := magic.isValidImage(fhs[0]); !valid {
		return errors.New(fmt.Sprintf("Image must be between %s.", strings.Join(magic.AllowExt, ", ")))
	}

	// validation size image
	if fhs[0].Size > int64(magic.MaxFileSize) {
		return errors.New(fmt.Sprintf("An image cannot greater than %c Mb.", strconv.Itoa(magic.MaxFileSize)[0]))
	}

	// append to files
	magic.Files = fhs

	return nil
}

// validator for multiple image with addon feature all image must be unique.
func (magic *MagicImage) ValidateMultipleImage(key string) error {
	fhs := magic.MultipartForm.File[key]

	if !magic.Required && len(fhs) == 0 {
		return nil
	}

	if magic.Required && len(fhs) == 0 {
		return errors.New("Image is required.")
	}

	for index, value := range fhs {
		// validate images
		ext, valid := magic.isValidImage(value)
		if !valid {
			return errors.New(fmt.Sprintf("The image at index %d must be between %s.", index+1, strings.Join(magic.AllowExt, ", ")))
		}

		// validation size image
		if value.Size > int64(magic.MaxFileSize) {
			return errors.New(fmt.Sprintf("An image at index %d cannot greater than %c Mb.", index+1, strconv.Itoa(magic.MaxFileSize)[0]))
		}

		magic.appendImageHash(ext, value)
	}

	// image must be unique in a list of images
	if magic.findDuplicateImage() {
		return errors.New("Each image must be unique.")
	}

	// check minimum or maximum image in list
	if magic.MinFileInSlice != 0 && len(fhs) < magic.MinFileInSlice {
		return errors.New(fmt.Sprintf("At least %d image must be upload.", magic.MinFileInSlice))
	}

	if magic.MaxFileInSlice != 0 && len(fhs) > magic.MaxFileInSlice {
		return errors.New(fmt.Sprintf("Maximal %d images to be upload.", magic.MaxFileInSlice))
	}

	// append to files
	magic.Files = fhs

	return nil
}

// reverseOrientation amply`s what ever operation is necessary to transform given orientation
// to the orientation 1
func (magic *MagicImage) reverseOrientation(img image.Image, o string) *image.NRGBA {
	switch o {
	case "1":
		return imaging.Clone(img)
	case "2":
		return imaging.FlipV(img)
	case "3":
		return imaging.Rotate180(img)
	case "4":
		return imaging.Rotate180(imaging.FlipV(img))
	case "5":
		return imaging.Rotate270(imaging.FlipV(img))
	case "6":
		return imaging.Rotate270(img)
	case "7":
		return imaging.Rotate90(imaging.FlipV(img))
	case "8":
		return imaging.Rotate90(img)
	}
	return imaging.Clone(img)
}

func (magic *MagicImage) fixOrientation(path_upload, filename string, fhs *multipart.FileHeader) (string, error) {
	// check directory exists if not create
	if _, err := os.Stat(path_upload); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path_upload, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	filenameSlice := strings.Split(filename, ".")

	tempFile, err := ioutil.TempFile(path_upload, filenameSlice[0]+".*."+filenameSlice[1])
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	file, err := fhs.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	f, err := os.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer f.Close()

	// fix orientation
	x, err := exif.Decode(f)
	if err == nil {
		orientation, err := x.Get(exif.Orientation)
		if err == nil {
			// reverse orientation
			f, _ := imaging.Open(tempFile.Name())
			fo := magic.reverseOrientation(f, orientation.String())
			os.Remove(tempFile.Name())
			imaging.Save(fo, tempFile.Name())
		}
	}

	return tempFile.Name(), nil
}

func (magic *MagicImage) resizeImageAndSave(width, height int, filename string, square bool) {
	imgResize, _ := imaging.Open(filename)

	var imgResult *image.NRGBA

	if square {
		imgResult = imaging.Fill(imgResize, width, height, imaging.Center, imaging.Lanczos)
	} else {
		imgResult = imaging.Resize(imgResize, width, height, imaging.Lanczos)
	}

	os.Remove(filename)
	imaging.Save(imgResult, filename)
	// save filename
	SFilename := strings.Split(filename, "/")
	magic.FileNames = append(magic.FileNames, SFilename[len(SFilename)-1])
}

// save image for single image & multiple image with feature fix orientation iphone & create dir when doesn't exists.
func (magic *MagicImage) SaveImages(width, height int, path_upload string, square bool) error {
	for _, value := range magic.Files {
		ext, _ := magic.isValidImage(value)
		uuid := uuid.New()

		filename := uuid.String() + "." + ext

		fname, err := magic.fixOrientation(path_upload, filename, value)
		if err != nil {
			return err
		}
		magic.resizeImageAndSave(width, height, fname, square)
	}

	return nil
}
