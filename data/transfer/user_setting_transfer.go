package transfer

import (
	"time"
)

type UserSetting struct {
	Key         string    `json:"key"`
	Value       any       `json:"value"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserSettingUpdate struct {
	Value any `json:"new_setting_value" validate:"required"`
}
