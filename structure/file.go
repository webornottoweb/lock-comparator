package structure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/mcuadros/go-version"
)

// LockFile represents the whole composer.lock file
type LockFile struct {
	Readme           json.RawMessage `json:"_readme,omitempty"`
	Packages         []Package       `json:"packages,omitempty"`
	PackagesDev      []Package       `json:"packages-dev,omitempty"`
	ContentHash      string          `json:"content-hash,omitempty"`
	Aliases          json.RawMessage `json:"aliases,omitempty"`
	MinimumStability string          `json:"mininum-stability,omitempty"`
	StabilityFlags   json.RawMessage `json:"stability-flags,omitempty"`
	PreferStable     bool            `json:"prefer-stable,omitempty"`
	PreferLowest     bool            `json:"prefer-lowest,omitempty"`
	Platform         json.RawMessage `json:"platform,omitempty"`
	PlatformDev      json.RawMessage `json:"platform-dev,omitempty"`
}

// Merge received lock file with current file
// New file will be returned
func (lf *LockFile) Merge(lfRight *LockFile) *LockFile {
	var result LockFile

	result.Readme = lf.Readme
	result.Packages = lf.Packages
	result.PackagesDev = lf.PackagesDev
	result.ContentHash = lf.ContentHash
	result.Aliases = lf.Aliases
	result.MinimumStability = lf.MinimumStability
	result.StabilityFlags = lf.StabilityFlags
	result.PreferStable = lf.PreferStable
	result.PreferLowest = lf.PreferLowest
	result.Platform = lf.Platform
	result.PlatformDev = lf.PlatformDev

	for i := 0; i < len(lfRight.PackagesDev); i++ {
		foundDev, idxDev, errDev := result.getPackageDevByname(lfRight.PackagesDev[i].Name)
		found, idx, err := result.getPackageByname(lfRight.PackagesDev[i].Name)

		if err != nil && errDev != nil { // package was not found
			fmt.Println("Dev package was not found in composer_l.lock: " + lfRight.PackagesDev[i].Name)
			fmt.Println("Adding: " + lfRight.PackagesDev[i].Name)
			result.PackagesDev = append(result.PackagesDev, lfRight.PackagesDev[i])
		} else { // package was found and we'll show versions
			if err == nil { // package was found in prod section
				fmt.Println("Dev package was found in prod composer_l.lock: " + lfRight.PackagesDev[i].Name)

				comparison := version.CompareSimple(version.Normalize(found.Version), version.Normalize(lfRight.PackagesDev[i].Version))

				if comparison < 0 {
					fmt.Println("Adding into dev fresh: " + lfRight.PackagesDev[i].Name)
					result.PackagesDev = append(result.PackagesDev, lfRight.PackagesDev[i])
				} else {
					fmt.Println("Adding into dev old: " + lfRight.PackagesDev[i].Name)
					result.PackagesDev = append(result.PackagesDev, result.Packages[idx])
				}

				copy(result.Packages[idx:], result.Packages[idx+1:])
				result.Packages = result.Packages[:len(result.Packages)-1]
			} else { // package was found in dev section
				fmt.Println("Dev package was found in dev composer_l.lock: " + lfRight.PackagesDev[i].Name)
				comparison := version.CompareSimple(version.Normalize(foundDev.Version), version.Normalize(lfRight.PackagesDev[i].Version))
				if comparison < 0 {
					fmt.Println("Replace in dev: " + lfRight.PackagesDev[i].Name + "(" + version.Normalize(foundDev.Version) + "," + version.Normalize(lfRight.PackagesDev[i].Version) + ")")
					result.PackagesDev[idxDev] = lfRight.PackagesDev[i]
				}
			}
		}
	}

	for i := 0; i < len(lfRight.Packages); i++ {
		foundDev, idxDev, errDev := result.getPackageDevByname(lfRight.Packages[i].Name)
		found, idx, err := result.getPackageByname(lfRight.Packages[i].Name)

		if err != nil && errDev != nil { // package was not found
			fmt.Println("Package was not found in composer_l.lock: " + lfRight.Packages[i].Name)
			fmt.Println("Adding into prod: " + lfRight.Packages[i].Name)
			result.Packages = append(result.Packages, lfRight.Packages[i])
		} else { // package was found and we'll show versions
			if err != nil { // package was found in dev
				fmt.Println("Package was found in dev of composer_l.lock: " + lfRight.Packages[i].Name)
				comparison := version.CompareSimple(version.Normalize(foundDev.Version), version.Normalize(lfRight.Packages[i].Version))
				if comparison < 0 {
					fmt.Println("Replace in dev: " + lfRight.Packages[i].Name + "(" + version.Normalize(foundDev.Version) + ", " + version.Normalize(lfRight.Packages[i].Version) + ")")
					result.PackagesDev[idxDev] = lfRight.Packages[i]
				}
			} else {
				fmt.Println("Package was found in prod of composer_l.lock: " + lfRight.Packages[i].Name)
				comparison := version.CompareSimple(version.Normalize(found.Version), version.Normalize(lfRight.Packages[i].Version))
				if comparison < 0 {
					fmt.Println("Replace in prod: " + lfRight.Packages[i].Name + "(" + version.Normalize(found.Version) + "," + version.Normalize(lfRight.Packages[i].Version) + ")")
					result.Packages[idx] = lfRight.Packages[i]
				}
			}
		}
	}

	return &result
}

func (lf *LockFile) getPackageByname(name string) (Package, int, error) {
	for i := 0; i < len(lf.Packages); i++ {
		if lf.Packages[i].Name == name {
			return lf.Packages[i], i, nil
		}
	}

	var result Package
	return result, -1, errors.New("Package was not found")
}

func (lf *LockFile) getPackageDevByname(name string) (Package, int, error) {
	for i := 0; i < len(lf.PackagesDev); i++ {
		if lf.PackagesDev[i].Name == name {
			return lf.PackagesDev[i], i, nil
		}
	}

	var result Package
	return result, -1, errors.New("Package was not found")
}

// Override composer file struct parsing into string
func (lf *LockFile) String() string {
	var b bytes.Buffer
	writer := io.Writer(&b)

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.Encode(lf)

	return b.String()
}
