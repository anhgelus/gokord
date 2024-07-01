package gokord

import (
	"encoding/json"
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"slices"
	"strconv"
	"strings"
)

type Version struct {
	Major      uint
	Minor      uint
	Patch      uint
	PreRelease string
}

type InnovationJson struct {
	Version  string              `json:"version"`
	Commands *InnovationCommands `json:"commands"`
}

type Innovation struct {
	Version  *Version
	Commands *InnovationCommands
}

type InnovationCommands struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Updated []string `json:"updated"`
}

var NilVersion = Version{Major: 0, Minor: 0, Patch: 0}

// LoadInnovationFromJson provided (could be embedded with go/embed)
func LoadInnovationFromJson(b []byte) ([]*Innovation, error) {
	var i []*Innovation
	var j []*InnovationJson
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}
	for _, item := range j {
		v, err := ParseVersion(item.Version)
		if err != nil {
			return nil, err
		}
		i = append(i, &Innovation{
			Version:  &v,
			Commands: item.Commands,
		})
	}
	return i, nil
}

func GetCommandsUpdate(bot *Bot) *InnovationCommands {
	remaining := bot.Innovations
	slices.SortFunc(remaining, func(a, b *Innovation) int {
		return a.Version.ForSort(b.Version)
	})
	slices.Reverse(remaining)
	utils.SendDebug("list of innovations", "innovations", remaining)
	l := len(remaining)
	if l == 0 {
		utils.SendSuccess("No updates available")
		return &InnovationCommands{}
	}
	lat := remaining[0]
	if lat == nil || lat.Version == nil {
		return nil
	}
	// loading bot data
	botData := BotData{Name: bot.Name}
	err := botData.Load()
	if err != nil {
		utils.SendAlert("version.go - Loading bot data for commands update", err.Error(), "name", bot.Name)
		return nil
	}
	// parse version of the bot
	ver, err := ParseVersion(botData.Version)
	if err != nil {
		utils.SendAlert(
			"version.go - Parsing version to compare for commands update",
			err.Error(),
			"version",
			botData.Version,
		)
		return nil
	}
	// if there is no update to do
	utils.SendDebug("last version and version of bot", "last", lat, "version of bot", ver)
	if !ver.Is(&NilVersion) {
		if lat.Version.Is(&ver) {
			utils.SendSuccess("No updates available")
			return &InnovationCommands{}
		} else if !lat.Version.NewerThan(&ver) {
			utils.SendWarn(
				"Bot has a newer version than the latest version available",
				"bot_version",
				botData.Version,
				"innovation_version",
				lat.Version,
			)
			return &InnovationCommands{}
		}
	}
	// get available versions
	var after []*Innovation
	version := lat
	id := 0
	for version.Version.NewerThan(&ver) {
		after = append(after, version)
		id++
		if id == len(remaining) {
			break
		}
		version = remaining[id]
	}
	// generate innovation commands
	slices.Reverse(after)
	cmds := &InnovationCommands{
		Added:   []string{},
		Removed: []string{},
		Updated: []string{},
	}
	for _, i := range after {
		for _, cmd := range i.Commands.Added {
			if slices.Contains(cmds.Removed, cmd) {
				// remove from "removed" and add to "updated"
				id = slices.Index(cmds.Removed, cmd)
				slices.Replace(cmds.Removed, id, id+1)
				cmds.Updated = append(cmds.Updated, cmd)
			} else {
				cmds.Added = append(cmds.Added, cmd)
			}
		}
		for _, cmd := range i.Commands.Updated {
			if slices.Contains(cmds.Removed, cmd) {
				// remove from "removed" and add to "added"
				id = slices.Index(cmds.Removed, cmd)
				slices.Replace(cmds.Removed, id, id+1)
				cmds.Added = append(cmds.Added, cmd)
			} else if slices.Contains(cmds.Added, cmd) {
				// do nothing
			} else {
				cmds.Updated = append(cmds.Updated, cmd)
			}
		}
		for _, cmd := range i.Commands.Removed {
			if slices.Contains(cmds.Added, cmd) {
				// remove from "added"
				id = slices.Index(cmds.Added, cmd)
				slices.Replace(cmds.Added, id, id+1)
			} else if slices.Contains(cmds.Updated, cmd) {
				// remove from "updated" and add to "removed"
				id = slices.Index(cmds.Updated, cmd)
				slices.Replace(cmds.Updated, id, id+1)
				cmds.Removed = append(cmds.Removed, cmd)
			} else {
				cmds.Removed = append(cmds.Removed, cmd)
			}
		}
	}
	return cmds
}

func ParseVersion(s string) (Version, error) {
	// if given version string is empty
	if len(s) == 0 {
		return Version{Major: 0, Minor: 0, Patch: 0}, nil
	}
	spl := strings.Split(s, "-")
	var pre string
	if len(spl) >= 2 {
		pre = strings.Join(spl[1:], "-")
	}
	sp := strings.Split(spl[0], ".")
	ma, err := strconv.Atoi(sp[0])
	if err != nil {
		return Version{}, err
	}
	mi, err := strconv.Atoi(sp[1])
	if err != nil {
		return Version{}, err
	}
	pa, err := strconv.Atoi(sp[2])
	if err != nil {
		return Version{}, err
	}
	return Version{
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

func (v *Version) UpdateBotVersion(bot *Bot) {
	botData := BotData{Name: bot.Name}
	err := botData.Load()
	if err != nil {
		utils.SendAlert("version.go - Loading bot data for update version", err.Error(), "name", bot.Name)
		return
	}
	botData.Version = v.String()
	err = botData.Save()
	if err != nil {
		utils.SendAlert("version.go - Saving bot data for update version", err.Error(), "name", bot.Name)
		return
	}
	utils.SendSuccess("Updated version of the bot")
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

// ForSort returns:
//   - 0 if o and v are the same version
//   - 1 if v is newer than o
//   - -1 if o is newer than v
func (v *Version) ForSort(o *Version) int {
	if v.Is(o) {
		return 0
	} else if v.NewerThan(o) {
		return 1
	} else {
		return -1
	}
}
