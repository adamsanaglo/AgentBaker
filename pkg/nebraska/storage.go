package nebraska

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// StorageClient uploads COSI artifacts to Azure Blob Storage through an AFD endpoint.
// The endpoint must speak the Azure Blob REST API (AFD typically proxies transparently).
type StorageClient struct {
	afdUploadEndpoint string
	cred              azcore.TokenCredential
}

// NewStorageClient creates a new storage client.
// afdUploadEndpoint is the AFD upload endpoint (e.g., "https://storage-upload.example.com").
// The legacy (tokenProvider, httpClient) parameters are kept for callsite
// compatibility but ignored; uploads now use the azblob SDK with
// DefaultAzureCredential (or a credential supplied via NewStorageClientWithCred).
func NewStorageClient(afdUploadEndpoint string, _ TokenProvider, _ interface{}) *StorageClient {
	return &StorageClient{afdUploadEndpoint: afdUploadEndpoint}
}

// NewStorageClientWithCred creates a storage client with an explicit credential.
func NewStorageClientWithCred(afdUploadEndpoint string, cred azcore.TokenCredential) *StorageClient {
	return &StorageClient{afdUploadEndpoint: afdUploadEndpoint, cred: cred}
}

// UploadCOSIArtifact uploads a COSI artifact to blob storage through AFD using
// the Azure Blob SDK with chunked uploads, retries, and AAD authentication.
func (s *StorageClient) UploadCOSIArtifact(ctx context.Context, container, blobPath, artifactPath string) error {
	cred := s.cred
	if cred == nil {
		var err error
		cred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return fmt.Errorf("get default azure credential: %w", err)
		}
	}

	client, err := azblob.NewClient(s.afdUploadEndpoint, cred, nil)
	if err != nil {
		return fmt.Errorf("create azblob client for %s: %w", s.afdUploadEndpoint, err)
	}

	f, err := os.Open(artifactPath)
	if err != nil {
		return fmt.Errorf("open artifact %s: %w", artifactPath, err)
	}
	defer f.Close()

	contentType := "application/octet-stream"
	_ = contentType
	_, err = client.UploadFile(ctx, container, blobPath, f, &azblob.UploadFileOptions{
		BlockSize:   16 * 1024 * 1024,
		Concurrency: 8,
	})
	if err != nil {
		return fmt.Errorf("upload %s -> %s/%s/%s: %w", artifactPath, s.afdUploadEndpoint, container, blobPath, err)
	}
	return nil
}

// DeriveDownloadURL constructs the deterministic download URL for a COSI artifact.
// The download URL is: https://<afdDownloadHostname>/<container>/<blobPath>
func DeriveDownloadURL(afdDownloadHostname, container, blobPath string) string {
	return fmt.Sprintf("https://%s/%s/%s", afdDownloadHostname, container, blobPath)
}
