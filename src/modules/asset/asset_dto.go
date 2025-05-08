package asset

import (
	"time"

	"github.com/MetaDandy/cuent-ai-core/src/model"
	generatejob "github.com/MetaDandy/cuent-ai-core/src/modules/generate_job"
)

type AssetResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Video_URL  string `json:"video_url"`
	Audio_URL  string `json:"audio_url"`
	Line       string `json:"line"`
	AudioState string `json:"audio_state"`
	VideoState string `json:"video_state"`
	Duration   uint   `json:"duration"`
	Position   int    `json:"position"`

	Generated []generatejob.GeneratedJobResponse `json:"meta_data,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func AssetToDto(u *model.Asset) AssetResponse {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		t := u.DeletedAt.Time
		deletedAt = &t
	}

	var generated []generatejob.GeneratedJobResponse
	if len(u.GeneratedJobs) > 0 {
		generated = make([]generatejob.GeneratedJobResponse, 0, len(u.GeneratedJobs))
		generated = generatejob.GeneratedJobToLisDTO(u.GeneratedJobs)
	}

	return AssetResponse{
		ID:         u.ID.String(),
		Type:       u.Type,
		Video_URL:  u.Video_URL,
		Audio_URL:  u.Audio_URL,
		Line:       u.Line,
		AudioState: string(u.AudioState),
		VideoState: string(u.VideoState),
		Duration:   u.Duration,
		Position:   u.Position,

		Generated: generated,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func AssetsToListDTO(list []model.Asset) []AssetResponse {
	out := make([]AssetResponse, len(list))
	for i := range list {
		out[i] = AssetToDto(&list[i])
	}
	return out
}
