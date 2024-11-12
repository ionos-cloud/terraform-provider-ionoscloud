package ftp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/kardianos/ftps"
)

type Connection struct {
	Url               string // Server URL without any directory path. Example: ftp-fkb.ionos.com
	Port              int
	SkipVerify        bool           // Skip FTP server certificate verification. WARNING man-in-the-middle attack possible
	ServerCertificate *x509.CertPool // If FTP server uses self signed certificates, put this in tlsConfig. IONOS FTP Servers in prod DON'T need this
	Username          string
	Password          string
}

// IonosFtpUpload is a function that uploads an image to an IONOS FTP server.
// Contains IONOS-specific logic, i.e. the directory structure and naming conventions.
//
// image is the file to be uploaded.
// conn are the connection details to the FTP server.
func IonosFtpUpload(ctx context.Context, image *os.File, conn Connection) error {
	tlsConfig := tls.Config{
		InsecureSkipVerify: conn.SkipVerify,
		ServerName:         conn.Url,
		RootCAs:            conn.ServerCertificate,
		MaxVersion:         tls.VersionTLS12,
	}
	dialOptions := ftps.DialOptions{
		Host:        conn.Url,
		Port:        conn.Port,
		Username:    conn.Username,
		Passowrd:    conn.Password,
		ExplicitTLS: true,
		TLSConfig:   &tlsConfig,
	}

	c, err := ftps.Dial(ctx, dialOptions)
	if err != nil {
		return fmt.Errorf("dialing FTP server failed: %w", err)
	}

	// get path of image
	stat, err := image.Stat()
	if err != nil {
		return fmt.Errorf("failed while getting file info: %w", err)
	}

	// if extension is ISO or IMG, upload to /iso-images, else /hdd-images
	connPath := "hdd-images"
	if filepath.Ext(stat.Name()) == ".iso" || filepath.Ext(stat.Name()) == ".img" {
		connPath = "iso-images"
	}

	err = c.Chdir(connPath)
	if err != nil {
		return fmt.Errorf("failed while changing directory to %s within FTP server: %w", connPath, err)
	}

	files, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("failed while listing files within FTP server: %w", err)
	}

	// Check if there already exists an image with the given name at the location
	var errExists error
	for _, f := range files {
		if f.Name == stat.Name() {
			// Prepare an error - this MIGHT fail
			errExists = fmt.Errorf("%s might already exist at %s/%s. Please contact support at support@cloud.ionos.com to delete the old image - or choose a different image name. We're sorry for the inconvenience", stat.Name(), conn.Url, connPath)
		}
	}

	err = c.Upload(ctx, stat.Name(), image)
	if err != nil {
		err = fmt.Errorf("failed while uploading %s to FTP server: %w", stat.Name(), err)
		if errExists != nil {
			err = fmt.Errorf("%w\nNote: %w", err, errExists)
		}
		return err

	}

	return c.Close()
}

// CustomFtpUpload is a custom FTP upload function that can be used to upload files to a non-IONOS FTP server.
//
// image is the file to be uploaded.
// targetPath is the path where the file should be uploaded, including the desired filename.
// conn are the connection details to the FTP server.
func CustomFtpUpload(ctx context.Context, image *os.File, targetPath string, conn Connection) error {
	tlsConfig := tls.Config{
		InsecureSkipVerify: conn.SkipVerify,
		ServerName:         conn.Url,
		RootCAs:            conn.ServerCertificate,
		MaxVersion:         tls.VersionTLS12,
	}
	dialOptions := ftps.DialOptions{
		Host:        conn.Url,
		Port:        conn.Port,
		Username:    conn.Username,
		Passowrd:    conn.Password,
		ExplicitTLS: true,
		TLSConfig:   &tlsConfig,
	}

	c, err := ftps.Dial(ctx, dialOptions)
	if err != nil {
		return fmt.Errorf("dialing FTP server failed: %w", err)
	}

	err = c.Chdir(filepath.Dir(targetPath))
	if err != nil {
		return fmt.Errorf("failed while changing directory to %s within FTP server: %w", filepath.Dir(targetPath), err)
	}

	files, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("failed while listing files within FTP server: %w", err)
	}

	// Check if there already exists an image with the given name at the location
	desiredName := filepath.Base(targetPath)
	var errExists error
	for _, f := range files {
		if f.Name == desiredName {
			// Prepare an error - this MIGHT fail
			errExists = fmt.Errorf("%s might already exist at %s/%s. Please contact support at support@cloud.ionos.com to delete the old image - or choose a different image name. We're sorry for the inconvenience", desiredName, conn.Url, targetPath)
		}
	}

	err = c.Upload(ctx, desiredName, image)
	if err != nil {
		err = fmt.Errorf("failed while uploading %s to FTP server: %w", desiredName, err)
		if errExists != nil {
			err = fmt.Errorf("%w\nNote: %w", err, errExists)
		}
		return err

	}

	return c.Close()
}

// PollImage is a function that waits until an image is available at a given location, and returns the image's ID
func PollImage(ctx context.Context, meta any, imageName string, location string) (string, error) {
	client := meta.(services.SdkBundle).CloudApiClient

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("context cancelled or timed out: %w", ctx.Err())
		default:
			ls, _, err := client.ImagesApi.ImagesGet(ctx).
				Filter("public", "false").
				Filter("name", imageName).
				Filter("location", location).
				Execute()
			if err != nil {
				return "", fmt.Errorf("failed while fetching images: %w", err)
			}

			if len(*ls.Items) > 1 {
				// TODO: Why did FTP let us upload a duplicate image?
				panic(fmt.Errorf("multiple images with the same name found")) // TODO Exit gracefully
			}
			if len(*ls.Items) == 1 {
				return *(*ls.Items)[0].Id, nil
			}

			// Wait for 5 seconds before checking again
			time.Sleep(5 * time.Second)
		}
	}
}
var ValidFTPLocations = []string{"de/fra", "de/fkb", "de/txl", "us/las", "us/ewr", "es/vit"}
