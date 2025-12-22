package models

import (
	"time"

	"github.com/google/uuid"
)

type SkillCategory string

const (
	SkillCategoryBasics    SkillCategory = "basics"
	SkillCategoryFootwork  SkillCategory = "footwork"
	SkillCategoryBackward  SkillCategory = "backward"
	SkillCategoryWraps     SkillCategory = "wraps"
	SkillCategoryReleases  SkillCategory = "releases"
	SkillCategoryFloaters  SkillCategory = "floaters"
	SkillCategoryMultiples SkillCategory = "multiples"
)

func (sc SkillCategory) IsValid() bool {
	switch sc {
	case
		SkillCategoryBasics,
		SkillCategoryFootwork,
		SkillCategoryBackward,
		SkillCategoryWraps,
		SkillCategoryReleases,
		SkillCategoryFloaters,
		SkillCategoryMultiples:
		return true
	default:
		return false
	}
}

type Skill struct {
	ID uuid.UUID `json:"id"`

	Name  string `json:"name"`
	Level int16  `json:"level"`

	YoutubeVideoID   string `json:"youtubeVideoId"`
	IsVideoLandscape bool   `json:"isVideoLandscape"`

	Prerequisites []uuid.UUID     `json:"prerequisites"`
	Categories    []SkillCategory `json:"categories"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
