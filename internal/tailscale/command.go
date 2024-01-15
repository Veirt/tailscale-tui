package tailscale

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

var FinalCmd Command = Command{Name: "tailscale up", Flags: map[Flag]string{}}

type Flag struct {
	Name          string
	Desc          string
	IsBooleanFlag bool
}

func (f Flag) Title() string       { return f.Name }
func (f Flag) Description() string { return f.Desc }
func (f Flag) FilterValue() string { return f.Name }

type Command struct {
	Name  string
	Flags map[Flag]string
}

func (c Command) String() string {
	result := c.Name + " "

	// c.flags is a map of flag -> value
	// turn it into a string, sorted by flag name

	// sort the flags by name
	var keys []Flag
	for k := range c.Flags {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Name < keys[j].Name
	})

	for _, k := range keys {
		if k.IsBooleanFlag {
			result += k.Name + "=" + c.Flags[k] + " "
		} else {
			result += k.Name + " \"" + c.Flags[k] + "\" "
		}
	}

	return result
}

func CheckTailscale() error {
	_, err := exec.LookPath("tailscale")

	return err
}

func GetTailscaleUpFlags() []Flag {
	out, err := exec.Command("tailscale", "up", "--help").CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		if scanner.Text() == "FLAGS" {
			break
		}
	}

	cmds := []Flag{}
	for scanner.Scan() {
		out := scanner.Text()
		out = strings.TrimSpace(out)

		// if it starts with --, it's a flag
		if strings.HasPrefix(out, "--") {

			isBooleanFlag := strings.Contains(out, "false")
			name := ""
			if isBooleanFlag {
				name = out[:strings.IndexByte(out, ',')]
			} else {
				name = out[:strings.IndexByte(out, ' ')]
			}

			scanner.Scan()
			desc := strings.TrimSpace(scanner.Text())

			fl := Flag{Name: name, IsBooleanFlag: isBooleanFlag, Desc: desc}
			cmds = append(cmds, fl)

		}
	}

	return cmds

}
