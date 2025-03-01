package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/gui/keybindings"
	"github.com/jesseduffield/lazygit/pkg/theme"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type viewNameMapping struct {
	viewPtr **gocui.View
	name    string
}

func (gui *Gui) orderedViews() []*gocui.View {
	return lo.Map(gui.orderedViewNameMappings(), func(v viewNameMapping, _ int) *gocui.View {
		return *v.viewPtr
	})
}

func (gui *Gui) orderedViewNameMappings() []viewNameMapping {
	return []viewNameMapping{
		// first layer. Ordering within this layer does not matter because there are
		// no overlapping views
		{viewPtr: &gui.Views.Status, name: "status"},
		{viewPtr: &gui.Views.Snake, name: "snake"},
		{viewPtr: &gui.Views.Submodules, name: "submodules"},
		{viewPtr: &gui.Views.Files, name: "files"},
		{viewPtr: &gui.Views.Tags, name: "tags"},
		{viewPtr: &gui.Views.Remotes, name: "remotes"},
		{viewPtr: &gui.Views.Worktrees, name: "worktrees"},
		{viewPtr: &gui.Views.Branches, name: "localBranches"},
		{viewPtr: &gui.Views.RemoteBranches, name: "remoteBranches"},
		{viewPtr: &gui.Views.ReflogCommits, name: "reflogCommits"},
		{viewPtr: &gui.Views.Commits, name: "commits"},
		{viewPtr: &gui.Views.Stash, name: "stash"},
		{viewPtr: &gui.Views.SubCommits, name: "subCommits"},
		{viewPtr: &gui.Views.CommitFiles, name: "commitFiles"},

		{viewPtr: &gui.Views.Staging, name: "staging"},
		{viewPtr: &gui.Views.StagingSecondary, name: "stagingSecondary"},
		{viewPtr: &gui.Views.PatchBuilding, name: "patchBuilding"},
		{viewPtr: &gui.Views.PatchBuildingSecondary, name: "patchBuildingSecondary"},
		{viewPtr: &gui.Views.MergeConflicts, name: "mergeConflicts"},
		{viewPtr: &gui.Views.Secondary, name: "secondary"},
		{viewPtr: &gui.Views.Main, name: "main"},

		{viewPtr: &gui.Views.Extras, name: "extras"},

		// bottom line
		{viewPtr: &gui.Views.Options, name: "options"},
		{viewPtr: &gui.Views.AppStatus, name: "appStatus"},
		{viewPtr: &gui.Views.Information, name: "information"},
		{viewPtr: &gui.Views.Search, name: "search"},
		// this view takes up one character. Its only purpose is to show the slash when searching
		{viewPtr: &gui.Views.SearchPrefix, name: "searchPrefix"},

		// popups.
		{viewPtr: &gui.Views.CommitMessage, name: "commitMessage"},
		{viewPtr: &gui.Views.CommitDescription, name: "commitDescription"},
		{viewPtr: &gui.Views.Menu, name: "menu"},
		{viewPtr: &gui.Views.Suggestions, name: "suggestions"},
		{viewPtr: &gui.Views.Confirmation, name: "confirmation"},
		{viewPtr: &gui.Views.Tooltip, name: "tooltip"},

		// this guy will cover everything else when it appears
		{viewPtr: &gui.Views.Limit, name: "limit"},
	}
}

func (gui *Gui) createAllViews() error {
	frameRunes := []rune{'─', '│', '┌', '┐', '└', '┘'}
	switch gui.c.UserConfig.Gui.Border {
	case "double":
		frameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝'}
	case "rounded":
		frameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}
	case "hidden":
		frameRunes = []rune{' ', ' ', ' ', ' ', ' ', ' '}
	}

	var err error
	for _, mapping := range gui.orderedViewNameMappings() {
		*mapping.viewPtr, err = gui.prepareView(mapping.name)
		if err != nil && !gocui.IsUnknownView(err) {
			return err
		}
		(*mapping.viewPtr).FrameRunes = frameRunes
		(*mapping.viewPtr).FgColor = theme.GocuiDefaultTextColor
	}

	gui.Views.Options.FgColor = theme.OptionsColor
	gui.Views.Options.Frame = false

	gui.Views.SearchPrefix.BgColor = gocui.ColorDefault
	gui.Views.SearchPrefix.FgColor = gocui.ColorCyan
	gui.Views.SearchPrefix.Frame = false
	gui.c.SetViewContent(gui.Views.SearchPrefix, gui.Tr.SearchPrefix)

	gui.Views.Search.BgColor = gocui.ColorDefault
	gui.Views.Search.FgColor = gocui.ColorCyan
	gui.Views.Search.Editable = true
	gui.Views.Search.Frame = false
	gui.Views.Search.Editor = gocui.EditorFunc(gui.searchEditor)

	gui.Views.Stash.Title = gui.c.Tr.StashTitle

	gui.Views.Commits.Title = gui.c.Tr.CommitsTitle

	gui.Views.CommitFiles.Title = gui.c.Tr.CommitFiles

	gui.Views.Branches.Title = gui.c.Tr.BranchesTitle

	gui.Views.Remotes.Title = gui.c.Tr.RemotesTitle

	gui.Views.Worktrees.Title = gui.c.Tr.WorktreesTitle

	gui.Views.Tags.Title = gui.c.Tr.TagsTitle

	gui.Views.Files.Title = gui.c.Tr.FilesTitle

	for _, view := range []*gocui.View{gui.Views.Main, gui.Views.Secondary, gui.Views.Staging, gui.Views.StagingSecondary, gui.Views.PatchBuilding, gui.Views.PatchBuildingSecondary, gui.Views.MergeConflicts} {
		view.Title = gui.c.Tr.DiffTitle
		view.Wrap = true
		view.IgnoreCarriageReturns = true
		view.CanScrollPastBottom = gui.c.UserConfig.Gui.ScrollPastBottom
	}

	gui.Views.Staging.Title = gui.c.Tr.UnstagedChanges
	gui.Views.Staging.Highlight = false
	gui.Views.Staging.Wrap = true

	gui.Views.StagingSecondary.Title = gui.c.Tr.StagedChanges
	gui.Views.StagingSecondary.Highlight = false
	gui.Views.StagingSecondary.Wrap = true

	gui.Views.PatchBuilding.Title = gui.Tr.Patch
	gui.Views.PatchBuilding.Highlight = false
	gui.Views.PatchBuilding.Wrap = true

	gui.Views.PatchBuildingSecondary.Title = gui.Tr.CustomPatch
	gui.Views.PatchBuildingSecondary.Highlight = false
	gui.Views.PatchBuildingSecondary.Wrap = true

	gui.Views.MergeConflicts.Title = gui.c.Tr.MergeConflictsTitle
	gui.Views.MergeConflicts.Highlight = false
	gui.Views.MergeConflicts.Wrap = false

	gui.Views.Limit.Title = gui.c.Tr.NotEnoughSpace
	gui.Views.Limit.Wrap = true

	gui.Views.Status.Title = gui.c.Tr.StatusTitle

	gui.Views.AppStatus.BgColor = gocui.ColorDefault
	gui.Views.AppStatus.FgColor = gocui.ColorCyan
	gui.Views.AppStatus.Visible = false
	gui.Views.AppStatus.Frame = false

	gui.Views.CommitMessage.Visible = false
	gui.Views.CommitMessage.Title = gui.c.Tr.CommitSummary
	gui.Views.CommitMessage.Editable = true
	gui.Views.CommitMessage.Editor = gocui.EditorFunc(gui.commitMessageEditor)

	gui.Views.CommitDescription.Visible = false
	gui.Views.CommitDescription.Title = gui.c.Tr.CommitDescriptionTitle
	gui.Views.CommitDescription.Subtitle = utils.ResolvePlaceholderString(gui.Tr.CommitDescriptionSubTitle,
		map[string]string{
			"togglePanelKeyBinding": keybindings.Label(gui.UserConfig.Keybinding.Universal.TogglePanel),
		})
	gui.Views.CommitDescription.FgColor = theme.GocuiDefaultTextColor
	gui.Views.CommitDescription.Editable = true
	gui.Views.CommitDescription.Editor = gocui.EditorFunc(gui.commitDescriptionEditor)

	gui.Views.Confirmation.Visible = false
	gui.Views.Confirmation.Editor = gocui.EditorFunc(gui.promptEditor)

	gui.Views.Suggestions.Visible = false

	gui.Views.Menu.Visible = false

	gui.Views.Tooltip.Visible = false

	gui.Views.Information.BgColor = gocui.ColorDefault
	gui.Views.Information.FgColor = gocui.ColorGreen
	gui.Views.Information.Frame = false

	gui.Views.Extras.Title = gui.c.Tr.CommandLog
	gui.Views.Extras.Autoscroll = true
	gui.Views.Extras.Wrap = true

	gui.Views.Snake.Title = gui.c.Tr.SnakeTitle
	gui.Views.Snake.FgColor = gocui.ColorGreen

	return nil
}
