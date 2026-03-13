package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"regexp"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// UploadImage uploads a file to Cloudinary and returns the secure URL.
func UploadImage(ctx context.Context, file multipart.File) (string, error) {
	// cloudinary.New() automatically uses CLOUDINARY_URL env variable
	cld, err := cloudinary.New()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "gym-pro/foods",
		PublicID: uuid.New().String(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// DeleteImage deletes an image from Cloudinary using its URL or public ID.
func DeleteImage(ctx context.Context, secureURL string) error {
	// cloudinary.New() automatically uses CLOUDINARY_URL env variable
	cld, err := cloudinary.New()
	if err != nil {
		return fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	// Extract public_id from secure URL.
	// Example URL: https://res.cloudinary.com/<cloud_name>/image/upload/v123456789/gym-pro/foods/uuid-here.jpg
	// We need to extract "gym-pro/foods/uuid-here"
	publicID := extractPublicID(secureURL)
	if publicID == "" {
		return fmt.Errorf("could not extract public ID from URL: %s", secureURL)
	}

	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}

	return nil
}

// extractPublicID tries to extract the Cloudinary public ID from a standard secure URL.
// This is a naive implementation; assuming the folder structure is known.
func extractPublicID(secureURL string) string {
	// A more robust implementation would parse the URL properly.
	importRegex := `gym-pro/foods/[^.]+`
	
	re := regexp.MustCompile(importRegex)
	match := re.FindString(secureURL)
	return match
}
