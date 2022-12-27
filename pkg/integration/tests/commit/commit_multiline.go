package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitMultiline = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with a multi-line commit message",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.Model().CommitCount(0)

		input.PrimaryAction()
		input.Press(keys.Files.CommitChanges)

		input.CommitMessagePanel().Type("first line").AddNewline().AddNewline().Type("third line").Confirm()

		assert.Model().CommitCount(1)
		assert.Model().HeadCommitMessage(Equals("first line"))

		input.SwitchToCommitsView()
		assert.Views().Main().Content(MatchesRegexp("first line\n\\s*\n\\s*third line"))
	},
})
