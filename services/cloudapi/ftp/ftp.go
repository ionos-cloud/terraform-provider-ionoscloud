package ftp

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/kardianos/ftps"
)

// UploadProperties contains info needed to initialize an FTP connection to IONOS server and upload an image.
type UploadProperties struct {
	ImageFileProperties
	ConnectionProperties
}

type ImageFileProperties struct {
	Path       string // File name with file extension included
	DataBuffer *bufio.Reader
}
type ConnectionProperties struct {
	Url               string // Server URL without any directory path. Example: ftp-fkb.ionos.com
	Port              int
	SkipVerify        bool           // Skip FTP server certificate verification. WARNING man-in-the-middle attack possible
	ServerCertificate *x509.CertPool // If FTP server uses self signed certificates, put this in tlsConfig. IONOS FTP Servers in prod DON'T need this
	Username          string
	Password          string
}

func FtpUpload(ctx context.Context, p UploadProperties) error {
	tlsConfig := tls.Config{
		InsecureSkipVerify: p.SkipVerify,
		ServerName:         p.Url,
		RootCAs:            p.ServerCertificate,
		MaxVersion:         tls.VersionTLS12,
	}
	dialOptions := ftps.DialOptions{
		Host:        p.Url,
		Port:        p.Port,
		Username:    p.Username,
		Passowrd:    p.Password,
		ExplicitTLS: true,
		TLSConfig:   &tlsConfig,
	}

	c, err := ftps.Dial(ctx, dialOptions)
	if err != nil {
		return fmt.Errorf("dialing FTP server failed. Check username & password. FTP server doesn't support usage of JWT token: %w", err)
	}

	err = c.Chdir(filepath.Dir(p.Path))
	if err != nil {
		return fmt.Errorf("failed while changing directory within FTP server: %w", err)
	}

	files, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("failed while listing files within FTP server: %w", err)
	}

	// Check if there already exists an image with the given name at the location
	desiredFileName := filepath.Base(p.Path)
	var errExists error
	for _, f := range files {
		if f.Name == desiredFileName {
			errExists = fmt.Errorf("%s might already exist at %s. Please contact support at support@cloud.ionos.com to delete the old image - or choose a different image name. We're sorry for the inconvenience", desiredFileName, p.Url)
		}
	}

	err = c.Upload(ctx, desiredFileName, p.DataBuffer)
	if err != nil {
		err = fmt.Errorf("failed while uploading %s to FTP server: %w", desiredFileName, err)
		if errExists != nil {
			err = fmt.Errorf("%w\nNote: %w", err, errExists)
		}
		return err

	}

	return c.Close()
}
