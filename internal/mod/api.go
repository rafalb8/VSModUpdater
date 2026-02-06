package mod

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/mod/semver"
)

var (
	ErrNoUpdate       = errors.New("no update")
	ErrNoModID        = errors.New("no modid")
	ErrInvalidSemVer  = errors.New("is not a valid Semantic Version")
	ErrPreReleaseSkip = errors.New("skipped pre-release version")
	ErrUnstableSkip   = errors.New("skipped pre-release game version")
)

type Response struct {
	Mod        Mod    `json:"mod"`
	StatusCode string `json:"statuscode,omitempty"`
}

type Releases struct {
	ReleaseID  int      `json:"releaseid,omitempty"`
	Mainfile   string   `json:"mainfile,omitempty"`
	Filename   string   `json:"filename,omitempty"`
	FileID     int      `json:"fileid,omitempty"`
	Downloads  int      `json:"downloads,omitempty"`
	Tags       []SemVer `json:"tags,omitempty"` // Supported game versions list
	ModIDStr   string   `json:"modidstr,omitempty"`
	ModVersion SemVer   `json:"modversion"`
	Created    string   `json:"created,omitempty"`
	Changelog  string   `json:"changelog,omitempty"`
}

type Mod struct {
	ModID           int        `json:"modid,omitempty"`
	AssetID         int        `json:"assetid,omitempty"`
	Name            string     `json:"name,omitempty"`
	Text            string     `json:"text,omitempty"`
	Author          string     `json:"author,omitempty"`
	UrlAlias        any        `json:"urlalias,omitempty"`
	LogoFilename    any        `json:"logofilename,omitempty"`
	LogoFile        any        `json:"logofile,omitempty"`
	LogoFileDB      any        `json:"logofiledb,omitempty"`
	HomepageUrl     any        `json:"homepageurl,omitempty"`
	SourcecodeUrl   any        `json:"sourcecodeurl,omitempty"`
	TrailervideoUrl any        `json:"trailervideourl,omitempty"`
	IssuetrackerUrl any        `json:"issuetrackerurl,omitempty"`
	WikiUrl         any        `json:"wikiurl,omitempty"`
	Downloads       int        `json:"downloads,omitempty"`
	Follows         int        `json:"follows,omitempty"`
	TrendingPoints  int        `json:"trendingpoints,omitempty"`
	Comments        int        `json:"comments,omitempty"`
	Side            string     `json:"side,omitempty"`
	Type            string     `json:"type,omitempty"`
	Created         string     `json:"created,omitempty"`
	LastReleased    string     `json:"lastreleased,omitempty"`
	LastModified    string     `json:"lastmodified,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	Releases        []Releases `json:"releases,omitempty"`
	Screenshots     []any      `json:"screenshots,omitempty"`
}

type SemVer struct {
	string
}

func NewSemVer(x string) (SemVer, error) {
	v := SemVer{x}
	v.Sanitize()

	if !v.IsValid() {
		return SemVer{}, fmt.Errorf("SemVer: '%s' %w", x, ErrInvalidSemVer)
	}
	return v, nil
}

func (v *SemVer) UnmarshalJSON(data []byte) error {
	v.string = string(data)
	v.string = strings.Trim(v.string, `"`)

	v.Sanitize()

	if !v.IsValid() {
		return fmt.Errorf("SemVer: '%s' %w", v.string, ErrInvalidSemVer)
	}
	return nil
}

func (v *SemVer) Sanitize() {
	if v != nil && v.string != "" && v.string[0] != 'v' {
		v.string = "v" + v.string
	}
}

func (v SemVer) IsValid() bool {
	return semver.IsValid(v.string)
}

func (v SemVer) Compare(x SemVer) int {
	return semver.Compare(v.string, x.string)
}

func (v SemVer) PreRelease() bool {
	return semver.Prerelease(v.string) != ""
}

func (v SemVer) String() string {
	v.Sanitize()
	return v.string
}

func GetLatestVersion(versions []SemVer) SemVer {
	if len(versions) == 0 {
		return SemVer{}
	}

	return slices.MaxFunc(versions, func(a, b SemVer) int { return a.Compare(b) })
}

func IsAllPreRelease(versions []SemVer) bool {
	if len(versions) == 0 {
		return false
	}

	for _, v := range versions {
		if !v.PreRelease() {
			return false
		}
	}
	return true
}
