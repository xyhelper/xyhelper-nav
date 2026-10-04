// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// NavCategory is the golang structure for table nav_category.
type NavCategory struct {
	Id          int       `json:"id"          orm:"id"          description:""` //
	Name        string    `json:"name"        orm:"name"        description:""` //
	Icon        string    `json:"icon"        orm:"icon"        description:""` //
	Description string    `json:"description" orm:"description" description:""` //
	Sort        int       `json:"sort"        orm:"sort"        description:""` //
	Status      int       `json:"status"      orm:"status"      description:""` //
	Visibility  int       `json:"visibility"  orm:"visibility"  description:""` //
	CreatedAt   time.Time `json:"createdAt"   orm:"created_at"  description:""` //
	UpdatedAt   time.Time `json:"updatedAt"   orm:"updated_at"  description:""` //
}
