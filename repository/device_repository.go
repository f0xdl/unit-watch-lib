package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/f0xdl/unit-watch-lib/domain"
	"gorm.io/gorm"
	"time"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(dbDial gorm.Dialector, automigrate bool) (*DeviceRepository, error) {
	db, err := gorm.Open(dbDial)
	if err != nil {
		return nil, err
	}
	if automigrate {
		err = db.AutoMigrate(&Device{}, &Group{}, &Owner{})
		if err != nil {
			return nil, err
		}
	}
	//TODO: postgres.Open(dsn)
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(5)
	//db.SetConnMaxLifetime(30 * time.Minute)
	return &DeviceRepository{db: db}, nil
}

func (r *DeviceRepository) GetStatus(ctx context.Context, uid string) (domain.DeviceStatus, error) {
	var d Device
	err := r.db.
		WithContext(ctx).
		First(&d, "uid = ?", uid).
		Error
	if err != nil {
		return 0, err
	}
	return domain.DeviceStatus(d.Status), nil
}

func (r *DeviceRepository) UpdateStatus(ctx context.Context, uid string, status domain.DeviceStatus) error {
	return r.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("status", status).Error
}

func (r *DeviceRepository) Get(ctx context.Context, uid string) (domain.Device, error) {
	var d Device
	err := r.db.WithContext(ctx).
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

func (r *DeviceRepository) GetChatIds(ctx context.Context, uid string) ([]int64, error) {
	var d Device
	err := r.db.WithContext(ctx).
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

func (r *DeviceRepository) GetDeviceChats(ctx context.Context, uid string) (map[int64]string, error) {
	var d Device
	err := r.db.WithContext(ctx).
		Preload("Groups").
		First(&d, "uid = ?", uid).
		Error
	if err != nil {
		return nil, err
	}
	result := map[int64]string{}
	for _, g := range d.Groups {
		result[g.ChatID] = g.Lang
	}
	return result, nil
}

func (r *DeviceRepository) UpdateOnline(ctx context.Context, uid string, online bool) error {
	return r.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("last_seen", time.Now()).
		Update("online", online).
		Error
}

func (r *DeviceRepository) CreateDevice(ctx context.Context, uid, hash string) error {
	return r.db.WithContext(ctx).
		Create(NewDevice(uid)).
		Error

}

func (r *DeviceRepository) SetActive(ctx context.Context, uid string, active bool) error {
	return r.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("active", active).
		Error
}

func (r *DeviceRepository) UpdateExpires(ctx context.Context, uid string, t time.Time) error {
	return r.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Update("expires_at", t).
		Error
}
func (r *DeviceRepository) UpdateInfo(ctx context.Context, uid, label, point string) error {
	return r.db.WithContext(ctx).
		Model(&Device{}).
		Where("uid = ?", uid).
		Updates(map[string]interface{}{
			"label": label,
			"point": point,
		}).Error
}

func (r *DeviceRepository) AssignGroups(ctx context.Context, uid string, chatIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var device Device
		err := tx.
			Where("uid = ?", uid).
			First(&device).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("device with UID %s not found", uid)
			}
			return fmt.Errorf("failed to find device: %w", err)
		}

		err = tx.
			Model(&device).
			Association("Groups").
			Clear()
		if err != nil {
			return fmt.Errorf("failed to clear existing groups: %w", err)
		}

		if len(chatIDs) == 0 {
			return nil
		}

		var groups []*Group
		for _, chatID := range chatIDs {
			var group Group
			err = tx.
				Where("chat_id = ?", chatID).
				First(&group).
				Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				group = Group{ChatID: chatID}
				err = tx.Create(&group).Error
				if err != nil {
					return fmt.Errorf("failed to create group with chatID %d: %w", chatID, err)
				}
			} else if err != nil {
				return fmt.Errorf("failed to find group with chatID %d: %w", chatID, err)
			}
			groups = append(groups, &group)
		}
		err = tx.
			Model(&device).
			Association("Groups").
			Append(groups)
		if err != nil {
			return fmt.Errorf("failed to assign groups to device: %w", err)
		}

		return nil
	})
}

func (r *DeviceRepository) CreateGroup(ctx context.Context, chatID int64) error {
	return r.db.WithContext(ctx).
		Create(&Group{ChatID: chatID}).
		Error
}

func (r *DeviceRepository) GetDevice(ctx context.Context, uid string) (*domain.Device, error) {
	var device Device
	err := r.db.WithContext(ctx).
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
