package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWiget     *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()

	a.Settings().SetTheme(&myTheme{})

	// create a window for the app
	win := a.NewWindow("Markdown")

	// get the UI
	edit, preview := cfg.makeUI()

	cfg.createMenuItems(win)

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// set window size
	win.Resize(fyne.Size{
		Width:  800,
		Height: 600,
	})
	win.CenterOnScreen()

	// show and run
	win.ShowAndRun()

}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWiget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown
	// edit.OnChanged = onChanged

	return edit, preview
}

/* func onChanged(str string) {
	cfg.PreviewWidget.ParseMarkdown(str)
} */

func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	saveAsMenuItem := fyne.NewMenuItem("Save as...", app.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			defer write.Close()
			write.Write([]byte(app.EditWiget.Text))
		}
	}
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {

			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				// user cancelled
				return
			}

			defer read.Close()

			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWiget.SetText(string(data))
			app.CurrentFile = read.URI()

			win.SetTitle(win.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				// user cancelled
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name your file with a .md extension!", win)
				return
			}

			// save file
			write.Write([]byte(app.EditWiget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false

		}, win)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()

	}
}
