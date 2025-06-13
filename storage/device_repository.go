package storage

import (
	"context"
	"github.com/f0xdl/unit-watch-lib/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type GormStorage struct {
	db *gorm.DB
}

func NewGormStorage(path string, cfg *gorm.Config, automigrate bool) (*GormStorage, error) {
	if cfg == nil {
		cfg = &gorm.Config{}
	}
	db, err := gorm.Open(sqlite.Open(path), cfg)
	if err != nil {
		return nil, err
	}
	if automigrate {
		err = db.AutoMigrate(&Device{}, &Group{}, &Owner{})
		if err != nil {
			return nil, err
		}
	}

	return &GormStorage{db: db}, nil
}

func (s *GormStorage) GetStatus(ctx context.Context, uid string) (domain.DeviceStatus, error) {
	var d Device
	err := s.db.
		WithContext(ctx).
		First(&d, "uid = ?", uid).
		Error
	if err != nil {
		return 0, err
	}
	return domain.DeviceStatus(d.Status), nil
}

func (s *GormStorage) UpdateStatus(ctx context.Context, uid string, status domain.DeviceStatus) error {
	return s.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("status", status).Error
}

func (s *GormStorage) Get(ctx context.Context, uid string) (domain.Device, error) {
	var d Device
	err := s.db.WithContext(ctx).
		First(&d, "uid = ?", uid).
		Error
	if err != nil {
		return domain.Device{}, err
	}
	return domain.Device{
		UID:      d.UID,
		Status:   domain.DeviceStatus(d.Status),
		LastSeen: d.LastSeen,
		Active:   d.Active,
		Point:    d.Point,
		Label:    d.Label,
		Online:   d.Online,
	}, nil
}

func (s *GormStorage) GetChatIds(ctx context.Context, uid string) ([]int64, error) {
	var d Device
	err := s.db.WithContext(ctx).
		Preload("Groups").
		First(&d, "uid = ?", uid).
		Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(d.Groups))
	for i, g := range d.Groups {
		ids[i] = g.ChatID
	}
	return ids, nil
}

func (s *GormStorage) UpdateOnline(ctx context.Context, uid string, online bool) error {
	return s.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("last_seen", time.Now()).
		Update("online", online).
		Error
}

func (s *GormStorage) CreateDevice(ctx context.Context, uid string) error {
	return s.db.WithContext(ctx).
		Create(NewDevice(uid)).
		Error

}

func (s *GormStorage) SetActive(ctx context.Context, uid string, active bool) error {
	return s.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("active", active).
		Error
}

func (s *GormStorage) UpdateExpires(ctx context.Context, uid string, t time.Time) error {
	return s.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("expires_at", t).
		Error
}
func (s *GormStorage) UpdateInfo(ctx context.Context, uid, label, point string) error {
	return s.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Updates(map[string]interface{}{
			"label": label,
			"point": point,
		}).Error
}

func (s *GormStorage) AssignGroups(ctx context.Context, uid string, chatIDs []int64) error {
	var device Device
	err := s.db.WithContext(ctx).
		First(&device, "uid = ?", uid).
		Error
	if err != nil {
		return err
	}

	var groups []*Group
	err = s.db.WithContext(ctx).
		Where("chat_id IN ?", chatIDs).
		Find(&groups).Error
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).
		Model(&device).
		Association("Groups").
		Replace(groups)
}

func (s *GormStorage) CreateGroup(ctx context.Context, chatID int64) error {
	return s.db.WithContext(ctx).
		Create(&Group{ChatID: chatID}).
		Error
}

func (s *GormStorage) GetDevice(ctx context.Context, uid string) (*domain.Device, error) {
	var device Device
	err := s.db.WithContext(ctx).
		First(&device, "uid = ?", uid).
		Error
	if err != nil {
		return nil, err
	}

	return &domain.Device{
		UID:       device.UID,
		Label:     device.Label,
		Active:    device.Active,
		Status:    domain.DeviceStatus(device.Status),
		OwnerId:   device.OwnerId,
		LastSeen:  device.LastSeen,
		Point:     device.Point,
		ExpiresAt: device.ExpiresAt,
		Online:    device.Online,
	}, nil
}
