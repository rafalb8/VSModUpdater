package mod

import "errors"

var (
	ErrNoUpdate = errors.New("no update")
	ErrNoModID  = errors.New("no modid")
)

type Response struct {
	Mod        Mod    `json:"mod,omitempty"`
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
	ModVersion string   `json:"modversion,omitempty"`
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
