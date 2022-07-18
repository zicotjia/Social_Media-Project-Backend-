package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
	"time"
)

func inititiateCloudinary() *cloudinary.Cloudinary {
	cld, _ := cloudinary.NewFromParams("dp6h9unkh", "563879417284723", "ZUg0U1wALRKczaE3aD28WIdt8WY")

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
