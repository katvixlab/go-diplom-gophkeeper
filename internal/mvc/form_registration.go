package mvc

import (
	"fmt"
	"net/mail"

	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/rivo/tview"
)

var (
	formRegistrationUser = tview.NewForm()
)

func createFormRegistrationUser(cu *UIController) {
	user := &pb.User{}
	var password string
	formRegistrationUser.AddInputField("Username", "", 40,
		nil,
		func(text string) { user.Username = text })
	formRegistrationUser.AddPasswordField("Password", "", 40, rune(42),
		func(text string) { user.Password = text })
	formRegistrationUser.AddPasswordField("Confirm password", "", 40, rune(42),
		func(text string) { password = text })

	formRegistrationUser.AddInputField("Email", "", 40,
		nil,
		func(text string) { user.Email = text })

	formRegistrationUser.AddButton("Save", func() {
		if validateUser(user, password) {
			user.Password = password
			err := cu.sn.Register(user)
			if err != nil {
				createModalError(err, PageRegistrationUser)
				return
			}
			cu.AddItemInfoList(fmt.Sprintf("The user: %s registered successful", user.Username))
			pagesMenu.SwitchToPage(PageMenu)
		}
	})

	formRegistrationUser.AddButton("Back", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})
	formRegistrationUser.SetBorder(true).SetTitle("Registration").SetTitleAlign(tview.AlignLeft)
}

func validateUser(user *pb.User, password string) bool {
	modalError.
		ClearButtons().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				pagesMenu.SwitchToPage(PageRegistrationUser)
			}
		}).SetTitle("Error")

	if user.Password != password {
		modalError.
			SetText("The passwords not equal")
		pagesMenu.SwitchToPage(PageError)
	} else if len(user.Password) < 8 {
		modalError.
			SetText("Password must be at least 8 characters")
		pagesMenu.SwitchToPage(PageError)
	} else if _, err := mail.ParseAddress(user.Email); err != nil {
		modalError.
			SetText("Email address not valid")
		pagesMenu.SwitchToPage(PageError)
	} else if len(user.Username) < 5 {
		modalError.
			SetText("Username must be at least 5 characters")
		pagesMenu.SwitchToPage(PageError)
	} else {
		return true
	}
	return false
}
