package db

import (
	"context"
	"time"

	"github.com/XiaoleC05/MusicBox/internal/model"
)

type CredentialsRepository struct{}

func NewCredentialsRepository() *CredentialsRepository {
	return &CredentialsRepository{}
}

func (r *CredentialsRepository) GetByUser(ctx context.Context, userID int64) (*model.Credentials, error) {
	query := `SELECT user_id, kugou_cookie, kugou_token, netease_cookie, qq_cookie, updated_at 
			  FROM musicbox.user_credentials 
			  WHERE user_id = $1`

	var c model.Credentials
	var updatedAt time.Time
	var kugouCookie, kugouToken, neteaseCookie, qqCookie []byte

	err := Pool.QueryRow(ctx, query, userID).Scan(
		&c.UserID, &kugouCookie, &kugouToken, &neteaseCookie, &qqCookie, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.KugouCookie = string(kugouCookie)
	c.KugouToken = string(kugouToken)
	c.NeteaseCookie = string(neteaseCookie)
	c.QQCookie = string(qqCookie)

	return &c, nil
}

func (r *CredentialsRepository) Upsert(ctx context.Context, userID int64, kugouCookie, kugouToken, neteaseCookie, qqCookie string) error {
	query := `INSERT INTO musicbox.user_credentials (user_id, kugou_cookie, kugou_token, netease_cookie, qq_cookie) 
			  VALUES ($1, $2, $3, $4, $5) 
			  ON CONFLICT (user_id) 
			  DO UPDATE SET 
				  kugou_cookie = COALESCE($2, musicbox.user_credentials.kugou_cookie),
				  kugou_token = COALESCE($3, musicbox.user_credentials.kugou_token),
				  netease_cookie = COALESCE($4, musicbox.user_credentials.netease_cookie),
				  qq_cookie = COALESCE($5, musicbox.user_credentials.qq_cookie),
				  updated_at = NOW()`

	_, err := Pool.Exec(ctx, query, userID, kugouCookie, kugouToken, neteaseCookie, qqCookie)
	return err
}

func (r *CredentialsRepository) GetStatus(ctx context.Context, userID int64) (*model.CredentialsStatus, error) {
	query := `SELECT kugou_cookie IS NOT NULL, netease_cookie IS NOT NULL, qq_cookie IS NOT NULL 
			  FROM musicbox.user_credentials 
			  WHERE user_id = $1`

	var status model.CredentialsStatus
	err := Pool.QueryRow(ctx, query, userID).Scan(&status.Kugou, &status.Netease, &status.QQ)
	if err != nil {
		return &model.CredentialsStatus{}, nil
	}

	return &status, nil
}
