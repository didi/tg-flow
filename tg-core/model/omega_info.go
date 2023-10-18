package model

type OmegaInfo struct {
	Bootstrap           string            `json:"boot_strap"`
	User                string            `json:"user"`
	Password            string            `json:"pass_word"`
	Topic               string            `json:"topic"`
	SceneId             int64             `json:"scene_id"`
	EventId             map[string]string `json:"event_id"`
	BuryData            map[string]string `json:"bury_data"`
	IsCollectUserAction bool              `json:"is_collect_user_action"`
}
