package mod

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

var (
	ErrNoUpdate      = errors.New("no update")
	ErrNoModID       = errors.New("no modid")
	ErrInvalidSemVer = errors.New("is not a valid Semantic Version")
	ErrPreReleaseSkip = errors.New("skipped pre-release version")
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
	Tags       []string `json:"tags,omitempty"`
	ModIDStr   string   `json:"modidstr,omitempty"`
	ModVersion SemVer   `json:"modversion,omitempty"`
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

type SemVer string

func SemVerFromString(x string) (SemVer, error) {
	if x != "" && x[0] != 'v' {
		x = "v" + x
	}
	
	if !semver.IsValid(x) {
		return "", fmt.Errorf("SemVer: '%s' %w", x, ErrInvalidSemVer)
	}
	
	return SemVer(x), nil
}

func (v *SemVer) UnmarshalJSON(data []byte) error {
	x := string(data)
	x = strings.Trim(x, `"`)

	if x != "" && x[0] != 'v' {
		x = "v" + x
	}

	if !semver.IsValid(x) {
		return fmt.Errorf("SemVer: '%s' %w", x, ErrInvalidSemVer)
	}

	*v = SemVer(x)
	return nil
}

func (v SemVer) Compare(x SemVer) int {
	return semver.Compare(string(v), string(x))
}

func (v SemVer) PreRelease() bool {
	return semver.Prerelease(string(v)) != ""
}

func (v SemVer) String() string {
	if v != "" && v[0] != 'v' {
		return "v" + string(v)
	}
	return string(v)
}
