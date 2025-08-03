package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	minimumFileSizeFlag = flag.Int64("min-size", 1024*1024, "Minimum file size to display (in bytes)") // default is 1 MB
)

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func add(target *tview.TreeNode, path string, minimumFileSize int64) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			panic(err)
		}

		if fileInfo.Size() < minimumFileSize {
			continue
		}

		var node *tview.TreeNode

		if file.IsDir() {
			node = tview.NewTreeNode(file.Name() + "/").
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			node.SetColor(tcell.ColorBlue)
		} else {
			size := fileInfo.Size()
			sizeStr := formatFileSize(size)
			node = tview.NewTreeNode(file.Name() + " (" + sizeStr + ")").
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			node.SetColor(tcell.ColorYellow)
		}
		target.AddChild(node)
	}
}

func main() {
	flag.Parse()

	minimumFileSize := int64(1024 * 1024) // 1 MB
	if *minimumFileSizeFlag > 0 {
		minimumFileSize = *minimumFileSizeFlag
	}

	app := tview.NewApplication()

	rootDir := "."
	root := tview.NewTreeNode(rootDir)

	tree := tview.NewTreeView()
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tree, 0, 1, true)

	refreshDirectory := func(dirNode *tview.TreeNode, dirPath string) {
		dirNode.ClearChildren()
		add(dirNode, dirPath, minimumFileSize)
	}

	showConfirmDialog := func(message string, onConfirm func()) {
		modal := tview.NewModal().
			SetText(message).
			AddButtons([]string{"Yes", "No"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					onConfirm()
				}
				app.SetRoot(flex, true).SetFocus(tree)
			})
		app.SetRoot(modal, false).SetFocus(modal)
	}

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}

		path := reference.(string)

		fileInfo, err := os.Stat(path)
		if err != nil {
			showConfirmDialog(fmt.Sprintf("Error accessing '%s': %v\nRefresh parent directory?", path, err), func() {
				parentPath := filepath.Dir(path)

				var refreshParentDir func(*tview.TreeNode, string) bool
				refreshParentDir = func(currentNode *tview.TreeNode, targetPath string) bool {
					if ref := currentNode.GetReference(); ref != nil {
						if ref.(string) == targetPath {
							refreshDirectory(currentNode, targetPath)
							return true
						}
					}

					for _, child := range currentNode.GetChildren() {
						if refreshParentDir(child, targetPath) {
							return true
						}
					}
					return false
				}

				if parentPath == "." {
					refreshDirectory(root, parentPath)
				} else {
					refreshParentDir(root, parentPath)
				}
			})
			return
		}

		if fileInfo.IsDir() {
			children := node.GetChildren()
			if len(children) == 0 {
				add(node, path, minimumFileSize)
			} else {
				node.SetExpanded(!node.IsExpanded())
			}
		} else {
			fileName := filepath.Base(path)
			showConfirmDialog(fmt.Sprintf("Delete file '%s'?", fileName), func() {
				err := os.Remove(path)
				if err != nil {
					showConfirmDialog(fmt.Sprintf("Error deleting file: %v", err), func() {})
				} else {
					parentPath := filepath.Dir(path)

					var refreshParentDir func(*tview.TreeNode, string) bool
					refreshParentDir = func(currentNode *tview.TreeNode, targetPath string) bool {
						if ref := currentNode.GetReference(); ref != nil {
							if ref.(string) == targetPath {
								refreshDirectory(currentNode, targetPath)
								return true
							}
						}

						for _, child := range currentNode.GetChildren() {
							if refreshParentDir(child, targetPath) {
								return true
							}
						}
						return false
					}

					if parentPath == "." {
						refreshDirectory(root, parentPath)
					} else {
						refreshParentDir(root, parentPath)
					}
				}
			})
		}
	})

	flex.SetBorder(true).
		SetTitle(" purgo - select files to delete, directories to expand/collapse ").
		SetTitleAlign(tview.AlignCenter)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		case tcell.KeyEscape:
			app.SetRoot(flex, true).SetFocus(tree)
			return nil
		}
		return event
	})

	add(root, rootDir, minimumFileSize)
	tree.SetRoot(root).SetCurrentNode(root)

	if err := app.SetRoot(flex, true).SetFocus(tree).Run(); err != nil {
		panic(err)
	}
}
