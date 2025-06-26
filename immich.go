package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type IMMICHClient struct {
	ImmichURL  string
	AlbumsKeys *AlbumsKeys
}

func NewIMMICHClient(immichURL string, albumsKeys *AlbumsKeys) *IMMICHClient {
	return &IMMICHClient{
		ImmichURL:  immichURL,
		AlbumsKeys: albumsKeys,
	}
}

func (c *IMMICHClient) request(endpoint, method, apiKey string, body interface{}, out interface{}) error {
	url := fmt.Sprintf("%s/api%s", c.ImmichURL, endpoint)
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// only start with albums add apiKey
	if method == http.MethodGet && strings.HasPrefix(endpoint, "/albums") {
		req.Header.Set("x-api-key", apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("immich api error on endpoint %s: %d %s: %s", url, resp.StatusCode, resp.Status, string(b))
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *IMMICHClient) GetAlbumInfo(albumID string, withoutAssets bool) (AlbumInfo, error) {
	var result AlbumInfo
	endpoint := fmt.Sprintf("/albums/%s?withoutAssets=%t", albumID, withoutAssets)
	apiKey := c.AlbumsKeys.GetAlbumKey(albumID)
	err := c.request(endpoint, http.MethodGet, apiKey, nil, &result)
	return result, err
}

func (c *IMMICHClient) GetSharedLinksInfo(key string) (SharedLinkInfo, error) {
	var result SharedLinkInfo
	endpoint := fmt.Sprintf("/shared-links/me?key=%s", key)
	err := c.request(endpoint, http.MethodGet, "", nil, &result)
	return result, err
}

func (c *IMMICHClient) GetAssetInfo(assetID string, shareKey string) (AssetInfo, error) {
	var result AssetInfo
	endpoint := fmt.Sprintf("/assets/%s?key=%s", assetID, shareKey)
	err := c.request(endpoint, http.MethodGet, "", nil, &result)
	return result, err
}

func (c *IMMICHClient) GetAssetThumbnail(assetID, size, shareKey string) ([]byte, error) {
	return c.GetAssetFile(
		fmt.Sprintf("/assets/%s/thumbnail", assetID),
		map[string]string{"size": size, "key": shareKey},
	)
}

func (c *IMMICHClient) GetAssetOriginal(assetID, shareKey string) ([]byte, error) {
	return c.GetAssetFile(
		fmt.Sprintf("/assets/%s/original", assetID),
		map[string]string{"key": shareKey},
	)
}

func (c *IMMICHClient) GetAssetFile(path string, query map[string]string) ([]byte, error) {
	endpoint := path
	if len(query) > 0 {
		q := url.Values{}
		for k, v := range query {
			q.Set(k, v)
		}
		endpoint = fmt.Sprintf("%s?%s", endpoint, q.Encode())
	}
	url := fmt.Sprintf("%s/api%s", c.ImmichURL, endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("immich api error on endpoint %s: %d %s: %s", url, resp.StatusCode, resp.Status, string(b))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return b, nil
}

type AlbumInfo struct {
	AlbumName                  string      `json:"albumName"`
	AlbumThumbnailAssetId      string      `json:"albumThumbnailAssetId"`
	AssetCount                 int         `json:"assetCount"`
	Assets                     []AssetInfo `json:"assets"`
	CreatedAt                  string      `json:"createdAt"`
	Description                string      `json:"description"`
	EndDate                    string      `json:"endDate"`
	HasSharedLink              bool        `json:"hasSharedLink"`
	ID                         string      `json:"id"`
	IsActivityEnabled          bool        `json:"isActivityEnabled"`
	LastModifiedAssetTimestamp string      `json:"lastModifiedAssetTimestamp"`
	Order                      string      `json:"order"`
	OwnerId                    string      `json:"ownerId"`
	Shared                     bool        `json:"shared"`
	StartDate                  string      `json:"startDate"`
	UpdatedAt                  string      `json:"updatedAt"`
}

type AssetInfo struct {
	ID               string    `json:"id"`
	DeviceAssetId    string    `json:"deviceAssetId"`
	OwnerId          string    `json:"ownerId"`
	DeviceId         string    `json:"deviceId"`
	LibraryId        string    `json:"libraryId"`
	Type             string    `json:"type"`
	OriginalPath     string    `json:"originalPath"`
	OriginalFileName string    `json:"originalFileName"`
	OriginalMimeType string    `json:"originalMimeType"`
	Thumbhash        string    `json:"thumbhash"`
	FileCreatedAt    string    `json:"fileCreatedAt"`
	FileModifiedAt   string    `json:"fileModifiedAt"`
	LocalDateTime    string    `json:"localDateTime"`
	UpdatedAt        string    `json:"updatedAt"`
	IsFavorite       bool      `json:"isFavorite"`
	IsArchived       bool      `json:"isArchived"`
	IsTrashed        bool      `json:"isTrashed"`
	Visibility       *string   `json:"visibility,omitempty"`
	Duration         string    `json:"duration"`
	ExifInfo         *ExifInfo `json:"exifInfo,omitempty"`
	LivePhotoVideoId *string   `json:"livePhotoVideoId,omitempty"`
	// People           *[]string `json:"people,omitempty"`
	Checksum    *string  `json:"checksum,omitempty"`
	IsOffline   *bool    `json:"isOffline,omitempty"`
	HasMetadata *bool    `json:"hasMetadata,omitempty"`
	DuplicateId *string  `json:"duplicateId,omitempty"`
	Resized     bool     `json:"resized"`
	Tags        []string `json:"tags,omitempty"`
}

type ExifInfo struct {
	Make                 *string  `json:"make,omitempty"`
	Model                *string  `json:"model,omitempty"`
	ExifImageWidth       *int     `json:"exifImageWidth,omitempty"`
	ExifImageHeight      *int     `json:"exifImageHeight,omitempty"`
	FileSizeInByte       *int     `json:"fileSizeInByte,omitempty"`
	Orientation          *string  `json:"orientation,omitempty"`
	DateTimeOriginal     *string  `json:"dateTimeOriginal,omitempty"`
	ModifyDate           *string  `json:"modifyDate,omitempty"`
	TimeZone             *string  `json:"timeZone,omitempty"`
	LensMake             *string  `json:"lensMake,omitempty"`
	LensModel            *string  `json:"lensModel,omitempty"`
	FNumber              *float64 `json:"fNumber,omitempty"`
	FocalLength          *float64 `json:"focalLength,omitempty"`
	ISO                  *int     `json:"iso,omitempty"`
	ExposureTime         *string  `json:"exposureTime,omitempty"`
	ExposureCompensation *string  `json:"exposureCompensation,omitempty"`
	Latitude             *float64 `json:"latitude,omitempty"`
	Longitude            *float64 `json:"longitude,omitempty"`
	City                 *string  `json:"city,omitempty"`
	State                *string  `json:"state,omitempty"`
	Country              *string  `json:"country,omitempty"`
	Description          *string  `json:"description,omitempty"`
	ProjectionType       *string  `json:"projectionType,omitempty"`
	Rating               *int     `json:"rating,omitempty"`
}

type SharedLinkInfo struct {
	Album         *AlbumInfo  `json:"album,omitempty"`
	Assets        []AssetInfo `json:"assets,omitempty"`
	CreatedAt     *string     `json:"createdAt,omitempty"`
	ExpiresAt     *string     `json:"expiresAt,omitempty"`
	ID            *string     `json:"id,omitempty"`
	Key           *string     `json:"key,omitempty"`
	AllowDownload *bool       `json:"allowDownload,omitempty"`
	AllowUpload   *bool       `json:"allowUpload,omitempty"`
}
