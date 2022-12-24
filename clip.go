package main

import (
	"github.com/google/uuid"
)

type clip struct {
	FullScreen          bool
	Icon                []byte `plist:",omitempty" json:",omitempty"`
	IgnoreManifestScope bool
	IsRemovable         bool
	Label               string
	PayloadDescription  string // Configures settings for a web clip
	PayloadDisplayName  string // Web Clip
	PayloadIdentifier   string // com.apple.webClip.managed.uuid
	PayloadType         string // com.apple.webClip.managed
	PayloadUUID         uuid.UUID
	PayloadVersion      uint // 1
	Precomposed         bool
	URL                 string
}

func NewClip(label string, icon []byte, url string) *clip {
	c := &clip{
		FullScreen:          true,
		Icon:                icon,
		IgnoreManifestScope: true,
		IsRemovable:         true,
		Label:               label,
		PayloadDescription:  "Configures settings for a web clip",
		PayloadDisplayName:  "Web Clip",
		PayloadIdentifier:   "com.apple.webClip.managed." + uuid.NewString(),
		PayloadType:         "com.apple.webClip.managed",
		PayloadUUID:         uuid.New(),
		PayloadVersion:      1,
		Precomposed:         false,
		URL:                 url,
	}

	return c
}
