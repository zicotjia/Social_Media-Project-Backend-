package cloudinary

import "github.com/cloudinary/cloudinary-go/v2"

func inititiateCloudinary() *cloudinary.Cloudinary {
	cld, _ := cloudinary.NewFromParams("dp6h9unkh", "563879417284723", "ZUg0U1wALRKczaE3aD28WIdt8WY")

	return cld
}

var Cloudi *cloudinary.Cloudinary = inititiateCloudinary()
