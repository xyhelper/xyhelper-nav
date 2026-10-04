// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// NavSite is the golang structure for table nav_site.
type NavSite struct {
	Id            int       `json:"id"            orm:"id"              description:""` //
	CategoryId    int       `json:"categoryId"    orm:"category_id"     description:""` //
	Name          string    `json:"name"          orm:"name"            description:""` //
	Url           string    `json:"url"           orm:"url"             description:""` //
	Icon          string    `json:"icon"          orm:"icon"            description:""` //
	Description   string    `json:"description"   orm:"description"     description:""` //
	Sort          int       `json:"sort"          orm:"sort"            description:""` //
	Status        int       `json:"status"        orm:"status"          description:""` //
	Hits          int       `json:"hits"          orm:"hits"            description:""` //
	IconKey       string    `json:"iconKey"       orm:"icon_key"        description:""` //
	IconFetchedAt time.Time `json:"iconFetchedAt" orm:"icon_fetched_at" description:""` //
	CreatedAt     time.Time `json:"createdAt"     orm:"created_at"      description:""` //
	UpdatedAt     time.Time `json:"updatedAt"     orm:"updated_at"      description:""` //
}
