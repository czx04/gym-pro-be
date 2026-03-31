package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"regexp"
	"strings"

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

func UploadAvatarImage(ctx context.Context, file multipart.File, userID uuid.UUID) (string, error) {
	// cloudinary.New() automatically uses CLOUDINARY_URL env variable
	cld, err := cloudinary.New()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   fmt.Sprintf("gym-pro/avatars/%s", userID.String()),
		PublicID: fmt.Sprintf("%s_%s", userID.String(), uuid.New().String()),
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
func extractPublicID(secureURL string) string {
	secureURL = strings.TrimSpace(secureURL)
	if secureURL == "" {
		return ""
	}

	// If caller passed a public_id instead of a URL, accept it.
	if !strings.Contains(secureURL, "://") {
		return strings.TrimPrefix(secureURL, "/")
	}

	u, err := url.Parse(secureURL)
	if err != nil {
		return ""
	}

	// Cloudinary URL path example:
	// /<resource_type>/upload/v123456789/gym-pro/avatars/<userId>/<publicId>.jpg
	// We want gym-pro/avatars/<userId>/<publicId> (without extension).
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	uploadIdx := -1
	for i, p := range parts {
		if p == "upload" {
			uploadIdx = i
			break
		}
	}
	if uploadIdx == -1 || uploadIdx+1 >= len(parts) {
		return ""
	}

	after := parts[uploadIdx+1:]
	if len(after) == 0 {
		return ""
	}
	// Skip optional version segment (v123456...)
	if len(after) > 0 && regexp.MustCompile(`^v[0-9]+$`).MatchString(after[0]) {
		after = after[1:]
	}
	if len(after) == 0 {
		return ""
	}

	joined := strings.Join(after, "/")
	ext := path.Ext(joined)
	if ext != "" {
		joined = strings.TrimSuffix(joined, ext)
	}
	return joined
}
