package util

import (
	"encoding/base64"
	"errors"
	"github.com/HirbodBehnam/RedditDownloaderBot/config"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

// IsUrl checks if a string is an url
// From https://stackoverflow.com/a/55551215/4213397
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// FollowRedirect follows a page's redirect and returns the final URL
func FollowRedirect(u string) (string, error) {
	resp, err := config.GlobalHttpClient.Head(u)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	return resp.Request.URL.String(), nil
}

// DoesFfmpegExists returns true if ffmpeg is found
func DoesFfmpegExists() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// CheckFileSize checks the size of file before sending it to telegram
func CheckFileSize(f string, allowed int64) bool {
	fi, err := os.Stat(f)
	if err != nil {
		log.Println("Cannot get file size:", err.Error())
		return false
	}
	return fi.Size() <= allowed
}

// UUIDToBase64 uses the not standard base64 encoding to encode an uuid.UUID as string
// So instead of 36 chars we have 24
func UUIDToBase64(id uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(id[:])
}

// DownloadToFile downloads a link to a file
func DownloadToFile(link string, f *os.File) error {
	resp, err := config.GlobalHttpClient.Get(link)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = resp.Body.Close()
		return errors.New("forbidden")
	}
	_, err = io.Copy(f, resp.Body)
	_ = resp.Body.Close()
	return err
}
