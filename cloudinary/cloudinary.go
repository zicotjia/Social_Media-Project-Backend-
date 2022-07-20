package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func inititiateCloudinary() *cloudinary.Cloudinary {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUD"), os.Getenv("KEY"), os.Getenv("SECRET"))

	return cld
}

func CloudinaryUpload(file string, r chan string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	uploadParam, err := Cloudi.Upload.Upload(ctx, file, uploader.UploadParams{UseFilename: api.Bool(true)})
	if err != nil {
		log.Fatal(err)
	}
	r <- uploadParam.SecureURL
}

var Cloudi *cloudinary.Cloudinary = inititiateCloudinary()
