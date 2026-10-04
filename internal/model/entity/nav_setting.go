// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// NavSetting is the golang structure for table nav_setting.
type NavSetting struct {
	Id        int       `json:"id"        orm:"id"         description:""` //
	K         string    `json:"k"         orm:"k"          description:""` //
	V         string    `json:"v"         orm:"v"          description:""` //
	CreatedAt time.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt time.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
