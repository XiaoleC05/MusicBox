package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Song struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	Album    string   `json:"album"`
	Duration int      `json:"duration"`
	Platform string   `json:"platform"`
	PlatformID string `json:"platform_id"`
	Qualities []string `json:"qualities"`
}

type PlatformAdapter interface {
	Name() string
	Search(query string, page int) ([]Song, error)
	GetPlayURL(songID, quality string) (string, error)
	IsAvailable() bool
}

type KugouAdapter struct {
	client *http.Client
	cookie string
}

func NewKugouAdapter(cookie string) *KugouAdapter {
	return &KugouAdapter{
		client: &http.Client{Timeout: 10 * time.Second},
		cookie: cookie,
	}
}

func (k *KugouAdapter) Name() string {
	return "kugou"
}

func (k *KugouAdapter) IsAvailable() bool {
	return k.cookie != ""
}

func (k *KugouAdapter) Search(query string, page int) ([]Song, error) {
	if !k.IsAvailable() {
		return nil, fmt.Errorf("kugou adapter not configured")
	}

	apiURL := fmt.Sprintf("http://mobilecdn.kugou.com/api/v3/search/song?format=json&keyword=%s&page=%d&pagesize=20",
		url.QueryEscape(query), page)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", k.cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Status int `json:"status"`
		Data   struct {
			Info []struct {
				Hash     string `json:"hash"`
				SongName string `json:"songname"`
				Singer   string `json:"singername"`
				Album    string `json:"album_name"`
				Duration int    `json:"duration"`
			} `json:"info"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Status != 1 {
		return nil, fmt.Errorf("kugou API error: status=%d", result.Status)
	}

	var songs []Song
	for _, item := range result.Data.Info {
		songs = append(songs, Song{
			ID:         item.Hash,
			Title:      item.SongName,
			Artist:     item.Singer,
			Album:      item.Album,
			Duration:   item.Duration,
			Platform:   "kugou",
			PlatformID: item.Hash,
			Qualities:  []string{"standard", "high", "lossless"},
		})
	}

	return songs, nil
}

func (k *KugouAdapter) GetPlayURL(songID, quality string) (string, error) {
	if !k.IsAvailable() {
		return "", fmt.Errorf("kugou adapter not configured")
	}

	apiURL := fmt.Sprintf("http://mobilecdn.kugou.com/api/v3/song/info?hash=%s&format=json", songID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", k.cookie)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := k.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Status int `json:"status"`
		Data   struct {
			PlayURL string `json:"play_url"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Status != 1 {
		return "", fmt.Errorf("kugou API error: status=%d", result.Status)
	}

	if result.Data.PlayURL == "" {
		return "", fmt.Errorf("no play URL available")
	}

	return result.Data.PlayURL, nil
}
