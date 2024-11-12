package ftp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

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
		return fmt.Errorf("dialing FTP server failed. Check username & password. FTP server doesn't support usage of JWT token: %w", err)
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
