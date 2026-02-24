package mvc

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/rivo/tview"
)

var (
	formTextNote = tview.NewForm()
)

func createFormTextNote(cu *UIController, note models.TextNote) {
	formTextNote.Clear(true)
	var metaInfo string
	var textArea string
	formTextNote.AddTextArea("Text data", note.Text, 40, 0, 0,
		func(text string) { textArea = text })
	formTextNote.AddTextArea("Additional information", strings.Join(note.MetaInfo, "\n"), 40, 0, 0,
		func(text string) { metaInfo = text })
	formTextNote.AddInputField("Save as", note.NameRecord, 40,
		nil,
		func(text string) { note.NameRecord = text })

	formTextNote.AddButton("Save", func() {
		if note.Id == uuid.Nil {
			id, err := uuid.NewUUID()
			if err != nil {
				log.Fatal(err)
			}
			note.Id = id
		}
		if note.Created == 0 {
			note.Created = time.Now().Unix()
		}
		if textArea != "" {
			note.Text = textArea
		}
		if metaInfo != "" {
			note.MetaInfo = strings.Split(metaInfo, "\n")
		}

		note.Type = models.TEXT
		err := cu.AddNote(&note)
		if err != nil {
			createModalError(err, PageFormTextNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been saved with the text data:  %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})

	formTextNote.AddButton("Back", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})

	formTextNote.AddButton("Delete", func() {
		if note.Id == uuid.Nil {
			pagesMenu.SwitchToPage(PageMenu)
			return
		}
		err := cu.DeleteNote(note.Id)
		if err != nil {
			createModalError(err, PageFormTextNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been deleted:  %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})
	formTextNote.SetBorder(true).SetTitle("New text note").SetTitleAlign(tview.AlignLeft)
}
