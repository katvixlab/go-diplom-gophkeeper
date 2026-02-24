package mvc

import (
	"fmt"

	pb "github.com/katvixlab/go-diplom-gophkeeper/internal/interfaces/proto"
	"github.com/rivo/tview"
)

var (
	formAuthorization = tview.NewForm()
)

func createFormAuthorization(cu *UIController) {
	user := &pb.User{}
	formAuthorization.AddInputField("Email", "", 40,
		nil,
		func(text string) { user.Email = text })
	formAuthorization.AddPasswordField("Password", "", 40, rune(42),
		func(text string) { user.Password = text })

	formAuthorization.AddButton("Sign in", func() {
		err := cu.sn.Login(user)
		if err != nil {
			createModalError(err, PageSignIn)
			return
		}
		cu.AddItemInfoList(fmt.Sprintf("Welcome back %s!", user.Username))
		pagesMenu.SwitchToPage(PageMenu)

	})

	formAuthorization.AddButton("Cancel", func() {
		pagesMenu.SwitchToPage(PageMenu)
	})
	formAuthorization.SetBorder(true).SetTitle("Sig in").SetTitleAlign(tview.AlignLeft)
}
