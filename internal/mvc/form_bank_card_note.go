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
	formCardBankNote = tview.NewForm()
)

func createFormBankCardNote(cu *UIController, note models.BankCardNote) {
	formCardBankNote.Clear(true)
	var metaInfo string
	var cardNumber string
	formCardBankNote.AddInputField("Bank name", note.Bank, 40,
		nil,
		func(text string) { note.Bank = text })
	formCardBankNote.AddInputField("Card number", note.Number, 40,
		func(textToCheck string, lastChar rune) bool { return lastChar >= '0' && lastChar <= '9' },
		func(text string) { cardNumber = text })
	formCardBankNote.AddInputField("Expiration", note.Expiration, 40,
		nil,
		func(text string) { note.Expiration = text })
	formCardBankNote.AddInputField("Cardholder name", note.Cardholder, 40,
		nil,
		func(text string) { note.Cardholder = text })
	formCardBankNote.AddInputField("Security code", note.SecurityCode, 40,
		func(textToCheck string, lastChar rune) bool { return len(textToCheck) <= 3 },
		func(text string) { note.SecurityCode = text })
	formCardBankNote.AddTextArea("Additional information", strings.Join(note.MetaInfo, "\n"), 40, 0, 0,
		func(text string) { metaInfo = text })
	formCardBankNote.AddInputField("Save as", note.NameRecord, 40,
		nil,
		func(text string) { note.NameRecord = text })

	formCardBankNote.AddButton("Save", func() {
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
		if cardNumber != "" {
			note.Number = formatCardNumber(cardNumber)
		}

		note.Type = models.CARD
		err := cu.AddNote(&note)
		if err != nil {
			createModalError(err, PageFormBankCardNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been saved with the card data of the bank:%s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})

	formCardBankNote.AddButton("Back", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})

	formTextNote.AddButton("Delete", func() {
		if note.Id == uuid.Nil {
			pagesMenu.SwitchToPage(PageMenu)
			return
		}
		err := cu.DeleteNote(note.Id)
		if err != nil {
			createModalError(err, PageFormBankCardNote)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("The note has been deleted:  %s", note.NameRecord))
		pagesMenu.SwitchToPage(PageMenu)
	})
	formCardBankNote.SetBorder(true).SetTitle("New bank card note").SetTitleAlign(tview.AlignLeft)
}

func formatCardNumber(text string) string {
	str := ""
	i := 1
	for index, char := range text {
		str += string(char)
		if i%4 == 0 && index+1 != len(text) {
			str += "-"
		}
		i++
	}
	return str
}
