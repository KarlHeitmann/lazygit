package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitAndPush = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Make a commit and push",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		// No idea why I had to use version 2: it should be using my own computer's
		// font and the one iterm uses is version 3.
		config.UserConfig.Gui.NerdFontsVersion = "2"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("my-file.txt", "myfile content")
		shell.CreateFile("my-other-file.rb", "my-other-file content")

		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/demo")

		shell.CloneIntoRemote("origin")

		shell.SetBranchUpstream("feature/demo", "origin/feature/demo")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Stage a file")

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			SetCaptionPrefix("Commit our changes").
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("my commit summary").
			SwitchToDescription().
			Type("my commit description").
			SwitchToSummary().
			Confirm()

		t.Views().Commits().
			TopLines(
				Contains("my commit summary"),
			)

		t.SetCaptionPrefix("Push to the remote")

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Push)
	},
})
