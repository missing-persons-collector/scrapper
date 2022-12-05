package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
	"os"
	"time"
)

var cldInstance *cloudinary.Cloudinary

func NewCloudinary() error {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_CLOUD_SECRET"))

	if err != nil {
		return err
	}

	cldInstance = cld

	return nil
}

func Upload(file string, id string, country string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := cldInstance.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: id, AssetFolder: country, Folder: "croatia"})

	if err != nil {
		return err
	}

	return nil
}

func Exists(id string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	asset, err := cldInstance.Admin.Asset(ctx, admin.AssetParams{PublicID: "logo"})
	if err != nil {
		return "", err
	}

	return asset.SecureURL, nil
}

func RemoveImage(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ids := []string{id}

	_, err := cldInstance.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: ids})

	if err != nil {
		return err
	}

	return nil
}

func RemoveAllImages() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	_, err := cldInstance.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{Folder: "croatia"})

	if err != nil {
		log.Fatalln(err)
	}

}
