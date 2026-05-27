package repositories

import (
	"context"
	"mime/multipart"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
	"google.golang.org/api/iterator"
)

func (r *Repository) UploadImage(c fiber.Ctx, src multipart.File) (string, *fiber.Error) {
	result, err := r.Cloudinary.Cloudinary.Upload.Upload(c, src, uploader.UploadParams{Folder: "COMPRO NEED"})
	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return result.SecureURL, nil
}

// SEKARANG MENGGUNAKAN FIRESTORE!
func (r *Repository) InsertImageToFirestore(deviceValidID string, imageURL string) (string, *fiber.Error) {
	ctx := context.Background()

	// Kita simpan ke koleksi global "drainage_images" dengan auto-generated ID dokumen
	_, _, err := r.Firebase.Firestore.Collection("drainage_images").Add(ctx, map[string]interface{}{
		"device_id":    deviceValidID,
		"image_url":    imageURL,
		"created_at":   time.Now(),
		"last_updated": time.Now().Unix(),
	})

	if err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, "Gagal simpan metadata ke Firestore: "+err.Error())
	}

	return imageURL, nil
}

// AMBIL GAMBAR TERAKHIR DARI FIRESTORE (Query-nya jauh lebih simpel dibanding Flux!)
func (r *Repository) GetLatestImageFromFirestore() (map[string]interface{}, *fiber.Error) {
	ctx := context.Background()

	// Query: Urutkan berdasarkan created_at secara descending (terbaru di atas) dan batasi hanya 1 data
	query := r.Firebase.Firestore.Collection("drainage_images").
		OrderBy("created_at", firestore.Desc).
		Limit(1).
		Documents(ctx)

	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, fiber.NewError(fiber.StatusNotFound, "Belum ada data gambar yang tersimpan")
	}
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Gagal fetch data dari Firestore: "+err.Error())
	}

	return doc.Data(), nil
}
