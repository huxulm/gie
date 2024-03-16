package internal

import (
	"embed"
	"encoding/json"
	"regexp"

	_ "github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

//go:embed gitmoji.json
var fs embed.FS

type emoji struct {
	Emoji       string `json:"emoji"`
	Entity      string `json:"entity"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Semver      string `json:"semver"`
}

var emojix = map[string]*emoji{}
var emojiArr = []*emoji{}

func init() {
	r, _ := fs.Open("gitmoji.json")
	json.NewDecoder(r).Decode(&emojiArr)
	for i := range emojiArr {
		emojix[emojiArr[i].Code] = emojiArr[i]
	}
}

func New() *cobra.Command {
	git := cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			handleCommit(args)
			_ = GitExec(args...)
		},
	}
	/* git.AddCommand(&cobra.Command{
		DisableFlagParsing: true,
		Use:                "emoji",
		Aliases:            []string{"e"},
		Run: func(cmd *cobra.Command, args []string) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Code", "Emoji", "Description"})
			table.SetColMinWidth(2, 50)
			table.SetAutoWrapText(false)

			for i := range emojiArr {
				table.Append([]string{emojiArr[i].Code, emojiArr[i].Emoji, emojiArr[i].Description})
			}
			table.Render()
		},
	}) */
	git.DisableFlagParsing = true
	return &git
}

func handleCommit(args []string) {
	isCommit := false
	for i := range args {
		if args[i] == "commit" {
			isCommit = true
		}
		if isCommit && args[i] == "-m" && i+1 < len(args) {
			args[i+1] = emojify(args[i+1])
			break
		}
	}
}

var emojiReg = regexp.MustCompile(`:([a-z]|[A-Z]|[_])+:`)

func emojify(s string) string {
	return emojiReg.ReplaceAllStringFunc(s, func(t string) string {
		if emo, ok := emojix[t]; ok {
			return emo.Emoji
		}
		return t
	})
}
