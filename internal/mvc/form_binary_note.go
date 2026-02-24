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
	formBinaryNote = tview.NewForm()
)

func createFormBinaryNote(cu *UIController, note models.BinaryNote) {
	formBinaryNote.Clear(true)
	var metaInfo string
	var textArea string
	formBinaryNote.AddTextArea("Binary data", string(note.Binary), 40, 0, 0,
		func(text string) { textArea = text })
	formBinaryNote.AddTextArea("Additional information", strings.Join(note.MetaInfo, "\n"), 40, 0, 0,
		func(text string) { metaInfo = text })
	formBinaryNote.AddInputField("Save as", note.NameRecord, 40,
		nil,
		func(text string) { note.NameRecord = text })

	formBinaryNote.AddButton("Save", func() {
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
		if metaInfo != "" {
			note.MetaInfo = strings.Split(metaInfo, "\n")
		}
		if textArea != "" {
			note.Binary = []byte(textArea)
		}
		note.Type = models.BINARY
		err := cu.AddNote(&note)
		if err != nil {
			createModalError(err, PageFormBinaryNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been saved with the binary data: %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})

	formBinaryNote.AddButton("Back", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})

	formTextNote.AddButton("Delete", func() {
		if note.Id == uuid.Nil {
			pagesMenu.SwitchToPage(PageMenu)
			return
		}
		err := cu.DeleteNote(note.Id)
		if err != nil {
			createModalError(err, PageFormBinaryNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been deleted:  %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})
	formBinaryNote.SetBorder(true).SetTitle("New binary note").SetTitleAlign(tview.AlignLeft)
}
