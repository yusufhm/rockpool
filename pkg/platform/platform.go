package platform

import (
	"fmt"
	"runtime"
)

var (
	ConfigDir         string
	Name              string
	Domain            string
	Arch              = runtime.GOARCH
	UpgradeComponents []string
	NumTargets        int
	LagoonSshKey      string
)

func ToMap() map[string]string {
	return map[string]string{
		"Name":     Name,
		"Domain":   Domain,
		"Hostname": fmt.Sprintf("%s.%s", Name, Domain),
		"Arch":     Arch,
	}
}

func Hostname() string {
	return fmt.Sprintf("%s.%s", Name, Domain)
}

func TotalClusterNum() int {
	return NumTargets + 1
}
