package controllers

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"log"
)

//	var magicTable = map[string]string{
//		"\xff\xd8\xff":      "image/jpeg",
//		"\x89PNG\r\n\x1a\n": "image/png",
//		"GIF87a":            "image/gif",
//		"GIF89a":            "image/gif",
//	}
//
// // mimeFromIncipit returns the mime type of an image file from its first few
// // bytes or the empty string if the file does not look like a known file type
// func mimeFromIncipit(c *gin.Context) string {
//
//		incipt, _, _ := c.Request.FormFile("image")
//		defer incipt.Close()
//		buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
//		if _, err := incipt.Read(buff); err != nil {
//			fmt.Println(err) // do something with that error
//			c.Abort()
//		}
//
//		incipitStr := []byte(incipt)
//		for magic, mime := range magicTable {
//			if strings.HasPrefix(incipitStr, magic) {
//				return mime
//			}
//		}
//
//		return ""
//	}
func ImageResizing() *image.NRGBA {

	src, err := imaging.Open("me.jpg")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	src = imaging.CropAnchor(src, 400, 400, imaging.Center)

	src = imaging.Resize(src, 175, 0, imaging.Lanczos)

	img1 := imaging.AdjustContrast(src, 8)

	dst := imaging.New(300, 300, color.NRGBA{0, 0, 0, 0})

	dst = imaging.Paste(dst, img1, image.Pt(0, 0))

	// Save the resulting image as JPEG.
	err = imaging.Save(dst, "me.png")
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	return dst
}
