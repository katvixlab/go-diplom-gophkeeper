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
	formCredentialNote = tview.NewForm()
)

func createFormCredentialNote(cu *UIController, note models.CredentialNote) {
	formCredentialNote.Clear(true)
	var metaInfo string
	formCredentialNote.AddInputField("Username", note.Username, 40,
		nil,
		func(text string) { note.Username = text })
	formCredentialNote.AddInputField("Password", note.Password, 40,
		nil,
		func(text string) { note.Password = text })
	formCredentialNote.AddTextArea("Additional information", strings.Join(note.MetaInfo, "\n"), 40, 0, 0,
		func(text string) { metaInfo = text })
	formCredentialNote.AddInputField("Save as", note.NameRecord, 40,
		nil,
		func(text string) { note.NameRecord = text })

	formCredentialNote.AddButton("Save", func() {
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

		note.Type = models.CREDENTIAL
		err := cu.AddNote(&note)
		if err != nil {
			createModalError(err, PageFormCredential)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been saved with the cradential data: %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})

	formCredentialNote.AddButton("Back", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})

	formTextNote.AddButton("Delete", func() {
		if note.Id == uuid.Nil {
			pagesMenu.SwitchToPage(PageMenu)
			return
		}
		err := cu.DeleteNote(note.Id)
		if err != nil {
			createModalError(err, PageFormCredential)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been deleted:  %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})
	formCredentialNote.SetBorder(true).SetTitle("New credential note").SetTitleAlign(tview.AlignLeft)
}
