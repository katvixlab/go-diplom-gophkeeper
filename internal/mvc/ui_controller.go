package mvc

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/services/ui"
	"github.com/rivo/tview"
)

var (
	cu   *UIController
	once sync.Once
	log  *logger.Logger
	app  = tview.NewApplication()

	pagesMenu  = tview.NewPages()
	notesList  = tview.NewList().ShowSecondaryText(false)
	flexMain   = tview.NewFlex()
	modalError = tview.NewModal()
	textInfo   = tview.NewTextView()
)

type UIController struct {
	infoList []string
	sn       *ui.Service
}

const (
	PageMenu             = "Menu"
	PageFormBankCardNote = "Add Bank Card Note"
	PageFormCredential   = "Add Credential"
	PageFormTextNote     = "Add Text Note"
	PageFormBinaryNote   = "Add Binary Note"
	PageRegistrationUser = "Registration User"
	PageError            = "Error"
	PageSignIn           = "Sign in"
)

func NewUIController(logger *logger.Logger, serviceNote *ui.Service) *UIController {
	once.Do(func() {
		log = logger
		log.Infof("controller UI initializing")
		log.Infof("create controller UI")
		cu = &UIController{infoList: make([]string, 0), sn: serviceNote}
		log.Info("create menu")
		createMainMenu()
		log.Info("create flex")
		creteMainFlex()
		log.Info("setup input")
		setInput(cu)
	})
	return cu
}

func (cu *UIController) Run() error {
	log.Infof("controller UI running")
	return app.SetRoot(pagesMenu, true).EnableMouse(true).Run()
}

func (cu *UIController) AddItemInfoList(msg string) {
	cu.infoList = append(cu.infoList, msg)
	textInfo.Clear()
	textInfo.SetText(strings.Join(cu.infoList, "\n"))
	textInfo.SetTextColor(tcell.ColorYellowGreen).SetBorder(true).SetTitle("Info").SetBorderColor(tcell.ColorYellowGreen).SetTitleColor(tcell.ColorYellowGreen)
	textInfo.ScrollToEnd()
}

func (cu *UIController) AddNote(note models.Noteable) error {
	storage, err := cu.sn.AddNote(note)
	if err != nil {
		return err
	}
	createNotesList(*storage)
	return nil
}

func (cu *UIController) DeleteNote(id uuid.UUID) error {
	storage, err := cu.sn.DeleteNote(id)
	if err != nil {
		return err
	}
	createNotesList(*storage)
	return nil
}

func setInput(cu *UIController) *tview.Box {
	return flexMain.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 113:
			app.Stop()
		case 108:
			note, err := cu.sn.LoadNote()
			if err != nil {
				createModalError(err, PageMenu)
				return event
			}
			createNotesList(*note)
			cu.AddItemInfoList("Notes load is successful")
		case 98:
			formCardBankNote.Clear(true)
			createFormBankCardNote(cu, models.BankCardNote{})
			pagesMenu.SwitchToPage(PageFormBankCardNote)
		case 99:
			formCredentialNote.Clear(true)
			createFormCredentialNote(cu, models.CredentialNote{})
			pagesMenu.SwitchToPage(PageFormCredential)
		case 116:
			formTextNote.Clear(true)
			createFormTextNote(cu, models.TextNote{})
			pagesMenu.SwitchToPage(PageFormTextNote)
		case 105:
			formBinaryNote.Clear(true)
			createFormBinaryNote(cu, models.BinaryNote{})
			pagesMenu.SwitchToPage(PageFormBinaryNote)
		case 114:
			formRegistrationUser.Clear(true)
			createFormRegistrationUser(cu)
			pagesMenu.SwitchToPage(PageRegistrationUser)
		case 115:
			formAuthorization.Clear(true)
			createFormAuthorization(cu)
			pagesMenu.SwitchToPage(PageSignIn)
		}
		return event
	})
}

func createMainMenu() {
	pagesMenu.AddPage(PageMenu, flexMain, true, true)
	pagesMenu.AddPage(PageFormBankCardNote, createModalForm(formCardBankNote, 70, 23), true, false)
	pagesMenu.AddPage(PageFormCredential, createModalForm(formCredentialNote, 70, 17), true, false)
	pagesMenu.AddPage(PageFormTextNote, createModalForm(formTextNote, 70, 19), true, false)
	pagesMenu.AddPage(PageFormBinaryNote, createModalForm(formBinaryNote, 70, 19), true, false)
	pagesMenu.AddPage(PageRegistrationUser, createModalForm(formRegistrationUser, 70, 13), true, false)
	pagesMenu.AddPage(PageError, modalError, true, false)
	pagesMenu.AddPage(PageSignIn, createModalForm(formAuthorization, 55, 10), true, false)
}

func creteMainFlex() *tview.Flex {
	textMenu1 := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("(q) quit \n(l) load notes")
	textMenu2 := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("(b) add bank card \n(c) add credential")
	textMenu3 := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("(t) add text \n(i) add binary")
	textMenu4 := tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("(r) register an account \n(s) sign in")

	flexMain.SetBorder(true).SetTitle("Welcome to gopher keeper").SetTitleAlign(tview.AlignLeft)
	return flexMain.
		AddItem(notesList, 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(false).SetTitle(""), 0, 3, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(textMenu1, 0, 1, false).
				AddItem(textMenu2, 0, 1, false).
				AddItem(textMenu3, 0, 1, false).
				AddItem(textMenu4, 0, 1, false), 2, 1, false), 0, 2, false).
		AddItem(textInfo, 0, 1, false)
}

func createNotesList(storage []models.Noteable) {
	notesList.Clear()

	notesList.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		switch note := storage[i].(type) {
		case *models.BankCardNote:
			createFormBankCardNote(cu, *note)
			pagesMenu.SwitchToPage(PageFormBankCardNote)
		case *models.CredentialNote:
			createFormCredentialNote(cu, *note)
			pagesMenu.SwitchToPage(PageFormCredential)
		case *models.TextNote:
			createFormTextNote(cu, *note)
			pagesMenu.SwitchToPage(PageFormTextNote)
		case *models.BinaryNote:
			createFormBinaryNote(cu, *note)
			pagesMenu.SwitchToPage(PageFormBinaryNote)
		}
	})

	for i, note := range storage {
		item := fmt.Sprintf("[\r%s] %s", strings.ToUpper(note.GetType().String()), note.GetName())
		notesList.AddItem(item, "", rune(49+i), nil)
	}
}

func createModalForm(p tview.Primitive, width, height int) tview.Primitive {
	flex := tview.NewFlex()
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pagesMenu.SwitchToPage(PageMenu)
		}
		return event
	})
	return flex.
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)
}

func createModalError(err error, switchToPage string) {
	modalError.
		SetText(fmt.Sprintf("Error: %s", err.Error())).
		ClearButtons().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, buttonLabel string) {
			if buttonLabel == "OK" {
				pagesMenu.SwitchToPage(switchToPage)
			}
		}).SetTitle("Error")
	pagesMenu.SwitchToPage(PageError)
}
