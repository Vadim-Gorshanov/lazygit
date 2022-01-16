package gui

import (
	"log"

	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (gui *Gui) menuListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:        "menu",
			Key:             "menu",
			Kind:            types.PERSISTENT_POPUP,
			OnGetOptionsMap: gui.getMenuOptions,
		},
		GetItemsLength:      func() int { return gui.Views.Menu.LinesHeight() },
		OnGetPanelState:     func() types.IListPanelState { return gui.State.Panels.Menu },
		OnClickSelectedItem: gui.onMenuPress,
		Gui:                 gui,

		// no GetDisplayStrings field because we do a custom render on menu creation
	}
}

func (gui *Gui) filesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        FILES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return gui.State.FileManager.GetItemsLength() },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Files },
		OnFocus:         OnFocusWrapper(gui.onFocusFile),
		OnRenderToMain:  OnFocusWrapper(gui.filesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			lines := gui.State.FileManager.Render(gui.State.Modes.Diffing.Ref, gui.State.Submodules)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedFileNode()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        LOCAL_BRANCHES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Branches) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Branches },
		OnRenderToMain:  OnFocusWrapper(gui.branchesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetBranchListDisplayStrings(gui.State.Branches, gui.State.ScreenMode != SCREEN_NORMAL, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) remotesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:      func() int { return len(gui.State.Remotes) },
		OnGetPanelState:     func() types.IListPanelState { return gui.State.Panels.Remotes },
		OnRenderToMain:      OnFocusWrapper(gui.remotesRenderToMain),
		OnClickSelectedItem: gui.handleRemoteEnter,
		Gui:                 gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteListDisplayStrings(gui.State.Remotes, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedRemote()
			return item, item != nil
		},
	}
}

func (gui *Gui) remoteBranchesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        REMOTE_BRANCHES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.RemoteBranches) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.RemoteBranches },
		OnRenderToMain:  OnFocusWrapper(gui.remoteBranchesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetRemoteBranchListDisplayStrings(gui.State.RemoteBranches, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedRemoteBranch()
			return item, item != nil
		},
	}
}

func (gui *Gui) tagsListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        TAGS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Tags) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Tags },
		OnRenderToMain:  OnFocusWrapper(gui.tagsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetTagListDisplayStrings(gui.State.Tags, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedTag()
			return item, item != nil
		},
	}
}

func (gui *Gui) branchCommitsListContext() types.IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        BRANCH_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Commits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Commits },
		OnFocus:         OnFocusWrapper(gui.onCommitFocus),
		OnRenderToMain:  OnFocusWrapper(gui.branchCommitsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == BRANCH_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedLocalCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.Commits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
			)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedLocalCommit()
			return item, item != nil
		},
		RenderSelection: true,
	}
}

func (gui *Gui) subCommitsListContext() types.IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "branches",
			WindowName: "branches",
			Key:        SUB_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.SubCommits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.SubCommits },
		OnRenderToMain:  OnFocusWrapper(gui.subCommitsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			selectedCommitSha := ""
			if gui.currentContext().GetKey() == SUB_COMMITS_CONTEXT_KEY {
				selectedCommit := gui.getSelectedSubCommit()
				if selectedCommit != nil {
					selectedCommitSha = selectedCommit.Sha
				}
			}
			return presentation.GetCommitListDisplayStrings(
				gui.State.SubCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
				selectedCommitSha,
				startIdx,
				length,
				gui.shouldShowGraph(),
			)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedSubCommit()
			return item, item != nil
		},
		RenderSelection: true,
	}
}

func (gui *Gui) shouldShowGraph() bool {
	value := gui.UserConfig.Git.Log.ShowGraph
	switch value {
	case "always":
		return true
	case "never":
		return false
	case "when-maximised":
		return gui.State.ScreenMode != SCREEN_NORMAL
	}

	log.Fatalf("Unknown value for git.log.showGraph: %s. Expected one of: 'always', 'never', 'when-maximised'", value)
	return false
}

func (gui *Gui) reflogCommitsListContext() types.IListContext {
	parseEmoji := gui.UserConfig.Git.ParseEmoji
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commits",
			WindowName: "commits",
			Key:        REFLOG_COMMITS_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.FilteredReflogCommits) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.ReflogCommits },
		OnRenderToMain:  OnFocusWrapper(gui.reflogCommitsRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetReflogCommitListDisplayStrings(
				gui.State.FilteredReflogCommits,
				gui.State.ScreenMode != SCREEN_NORMAL,
				gui.cherryPickedCommitShaMap(),
				gui.State.Modes.Diffing.Ref,
				parseEmoji,
			)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedReflogCommit()
			return item, item != nil
		},
	}
}

func (gui *Gui) stashListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "stash",
			WindowName: "stash",
			Key:        STASH_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.StashEntries) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Stash },
		OnRenderToMain:  OnFocusWrapper(gui.stashRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetStashEntryListDisplayStrings(gui.State.StashEntries, gui.State.Modes.Diffing.Ref)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedStashEntry()
			return item, item != nil
		},
	}
}

func (gui *Gui) commitFilesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "commitFiles",
			WindowName: "commits",
			Key:        COMMIT_FILES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return gui.State.CommitFileManager.GetItemsLength() },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.CommitFiles },
		OnFocus:         OnFocusWrapper(gui.onCommitFileFocus),
		OnRenderToMain:  OnFocusWrapper(gui.commitFilesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			if gui.State.CommitFileManager.GetItemsLength() == 0 {
				return [][]string{{style.FgRed.Sprint("(none)")}}
			}

			lines := gui.State.CommitFileManager.Render(gui.State.Modes.Diffing.Ref, gui.Git.Patch.PatchManager)
			mappedLines := make([][]string, len(lines))
			for i, line := range lines {
				mappedLines[i] = []string{line}
			}

			return mappedLines
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedCommitFileNode()
			return item, item != nil
		},
	}
}

func (gui *Gui) submodulesListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "files",
			WindowName: "files",
			Key:        SUBMODULES_CONTEXT_KEY,
			Kind:       types.SIDE_CONTEXT,
		},
		GetItemsLength:  func() int { return len(gui.State.Submodules) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Submodules },
		OnRenderToMain:  OnFocusWrapper(gui.submodulesRenderToMain),
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSubmoduleListDisplayStrings(gui.State.Submodules)
		},
		SelectedItem: func() (types.ListItem, bool) {
			item := gui.getSelectedSubmodule()
			return item, item != nil
		},
	}
}

func (gui *Gui) suggestionsListContext() types.IListContext {
	return &ListContext{
		BasicContext: &BasicContext{
			ViewName:   "suggestions",
			WindowName: "suggestions",
			Key:        SUGGESTIONS_CONTEXT_KEY,
			Kind:       types.PERSISTENT_POPUP,
		},
		GetItemsLength:  func() int { return len(gui.State.Suggestions) },
		OnGetPanelState: func() types.IListPanelState { return gui.State.Panels.Suggestions },
		Gui:             gui,
		GetDisplayStrings: func(startIdx int, length int) [][]string {
			return presentation.GetSuggestionListDisplayStrings(gui.State.Suggestions)
		},
	}
}

func (gui *Gui) getListContexts() []types.IListContext {
	return []types.IListContext{
		gui.State.Contexts.Menu,
		gui.State.Contexts.Files,
		gui.State.Contexts.Branches,
		gui.State.Contexts.Remotes,
		gui.State.Contexts.RemoteBranches,
		gui.State.Contexts.Tags,
		gui.State.Contexts.BranchCommits,
		gui.State.Contexts.ReflogCommits,
		gui.State.Contexts.SubCommits,
		gui.State.Contexts.Stash,
		gui.State.Contexts.CommitFiles,
		gui.State.Contexts.Submodules,
		gui.State.Contexts.Suggestions,
	}
}
