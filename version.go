package gokord

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major      uint
	Minor      uint
	Patch      uint
	PreRelease string
}

type Innovation struct {
	Version  *Version              `json:"version"`
	Commands []*InnovationCommands `json:"commands"`
}

type InnovationCommands struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Updated []string `json:"updated"`
}

// LoadInnovationFromJson provided (could be embedded with go/embed)
func LoadInnovationFromJson(b []byte) ([]*Innovation, error) {
	var i []*Innovation
	err := json.Unmarshal(b, &i)
	return i, err
}

//func GetCommandsUpdate(bot *Bot) *InnovationCommands {
//	lat, id := getLatestInnovation(bot.Innovations)
//	if lat == nil {
//		return nil
//	}
//	botData := BotData{Name: bot.Name}
//	err := botData.Load()
//	if err != nil {
//		utils.SendAlert("version.go", "Loading bot data", "error", err.Error(), "name", bot.Name)
//		return nil
//	}
//	ver, err := ParseVersion(botData.Version)
//	if err != nil {
//		utils.SendAlert(
//			"version.go",
//			"Parsing version",
//			"error",
//			err.Error(),
//			"version",
//			botData.Version,
//		)
//		return nil
//	}
//	// if there is no update
//	if lat.Version.Is(ver) {
//		return &InnovationCommands{}
//	} else if !lat.Version.NewerThan(ver) {
//		utils.SendWarn(
//			"Bot has a newer version than the latest version in Innovation",
//			"bot_version",
//			botData.Version,
//			"innovation_version",
//			lat.Version,
//		)
//	}
//	var after []*Innovation
//	remaining := bot.Innovations
//	version := lat
//	for version.Version.NewerThan(ver) {
//		after = append(after, version)
//		var t []*Innovation
//		if id == len(remaining)-1 {
//			t = remaining[:id]
//		} else {
//			t = slices.Delete(remaining, id, id+1)
//		}
//		version, id = getLatestInnovation(t)
//	}
//	slices.Reverse(after)
//}

func getLatestInnovation(in []*Innovation) (*Innovation, int) {
	var lat *Innovation
	var id int
	for k := range len(in) {
		c := in[k]
		if k == 0 {
			lat = c
			id = k
		}
		if !lat.Version.NewerThan(c.Version) {
			lat = c
			id = k
		}
	}
	return lat, id
}

func ParseVersion(s string) (*Version, error) {
	sp := strings.Split(s, ".")
	ma, err := strconv.Atoi(sp[0])
	if err != nil {
		return nil, err
	}
	mi, err := strconv.Atoi(sp[1])
	if err != nil {
		return nil, err
	}
	pa, err := strconv.Atoi(sp[2])
	if err != nil {
		return nil, err
	}
	spl := strings.Split(s, "-")
	var pre string
	if len(spl) > 2 {
		pre = strings.Join(spl[1:], "-")
	}
	return &Version{
		Major:      uint(ma),
		Minor:      uint(mi),
		Patch:      uint(pa),
		PreRelease: pre,
	}, nil
}

func (v *Version) String() string {
	if len(v.PreRelease) != 0 {
		return fmt.Sprintf("%d.%d.%d-%s", v.Major, v.Minor, v.Patch, v.PreRelease)
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *Version) SetMajor(m uint) *Version {
	v.Major = m
	return v
}

func (v *Version) SetMinor(m uint) *Version {
	v.Minor = m
	return v
}

func (v *Version) SetPatch(p uint) *Version {
	v.Patch = p
	return v
}

func (v *Version) SetPreRelease(p string) *Version {
	v.PreRelease = p
	return v
}

// NewerThan check if the version is newer than version o
//
// Does not support pre-release checks
func (v *Version) NewerThan(o *Version) bool {
	if v.Major > o.Major {
		return true
	}
	if v.Minor > o.Minor {
		return true
	}
	if v.Patch > o.Patch {
		return true
	}
	return false
}

// Is check if this is the same version
func (v *Version) Is(o *Version) bool {
	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch && v.PreRelease == o.PreRelease
}
