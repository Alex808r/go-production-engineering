package srp_test

import (
	"context"
	"testing"

	"go-core-4/lessons/17-system-design/SOLID/1-SRP"
)

func TestCloudStorage_SRP(t *testing.T) {
	ctx := context.Background()

	var yandex srp.CloudStorage = srp.NewYandexDiskStorage("token-123")
	data, err := yandex.GetFile(ctx, "report.pdf")
	if err != nil {
		t.Fatalf("yandex GetFile failed: %v", err)
	}
	if string(data) != "yandex-data:report.pdf" {
		t.Errorf("unexpected yandex data: %s", string(data))
	}

	var google srp.CloudStorage = srp.NewGoogleDriveStorage("api-key-456")
	data, err = google.GetFile(ctx, "document.docx")
	if err != nil {
		t.Fatalf("google GetFile failed: %v", err)
	}
	if string(data) != "google-data:document.docx" {
		t.Errorf("unexpected google data: %s", string(data))
	}
}

func TestCloudStorage_Validation(t *testing.T) {
	ctx := context.Background()
	yandex := srp.NewYandexDiskStorage("token")

	_, err := yandex.GetFile(ctx, "")
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}

	err = yandex.UploadFile(ctx, "", []byte{1, 2, 3})
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}
