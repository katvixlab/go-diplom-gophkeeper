package mvc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/logger"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/models"
	"github.com/katvixlab/go-diplom-gophkeeper/internal/services/ui"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	baseNote = models.BaseNote{Id: uuid.Nil, NameRecord: "Test Note", Created: 1723652739, Type: models.CARD, MetaInfo: []string{"test", "test"}}
)

func TestNewUIController(t *testing.T) {
	type args struct {
		logger      *logger.Logger
		serviceNote *ui.Service
	}
	tests := []struct {
		name string
		args args
		want *UIController
	}{
		{
			name: "TestNewUIController",
			args: args{
				logger:      logger.NewLogger(logrus.New()),
				serviceNote: &ui.Service{},
			},
			want: &UIController{
				infoList: make([]string, 0),
				sn:       &ui.Service{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUIController(tt.args.logger, tt.args.serviceNote); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUIController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUIController_AddItemInfoList(t *testing.T) {
	type fields struct {
		infoList []string
		sn       *ui.Service
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestAddItemInfoList",
			fields: fields{
				infoList: make([]string, 0),
				sn:       &ui.Service{},
			},
			args: args{
				msg: "Hello World",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cu := &UIController{
				infoList: tt.fields.infoList,
				sn:       tt.fields.sn,
			}
			cu.AddItemInfoList(tt.args.msg)
		})
	}
}

func TestUIController_AddNote(t *testing.T) {
	conn, _ := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	service := ui.NewUIService(logger.NewLogger(logrus.New()), conn)
	type fields struct {
		infoList []string
		sn       *ui.Service
	}
	type args struct {
		note models.Noteable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestAddNote",
			fields: fields{
				infoList: make([]string, 0),
				sn:       service,
			},
			args: args{
				note: models.TextNote{
					Text:     "Hello World",
					BaseNote: baseNote,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cu := &UIController{
				infoList: tt.fields.infoList,
				sn:       tt.fields.sn,
			}
			if err := cu.AddNote(tt.args.note); (err != nil) != tt.wantErr {
				t.Errorf("AddNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createFormAuthorization(t *testing.T) {
	type args struct {
		cu *UIController
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateFormAuthorization",
			args: args{
				cu: &UIController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormAuthorization(tt.args.cu)
		})
	}
}

func Test_createFormBankCardNote(t *testing.T) {
	type args struct {
		cu   *UIController
		note models.BankCardNote
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateFormBankCardNote",
			args: args{
				cu: &UIController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormBankCardNote(tt.args.cu, tt.args.note)
		})
	}
}

func Test_createFormBinaryNote(t *testing.T) {
	type args struct {
		cu   *UIController
		note models.BinaryNote
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateFormBinaryNote",
			args: args{
				cu: &UIController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormBinaryNote(tt.args.cu, tt.args.note)
		})
	}
}

func Test_createFormCredentialNote(t *testing.T) {
	type args struct {
		cu   *UIController
		note models.CredentialNote
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateFormCredentialNote",
			args: args{
				cu: &UIController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormCredentialNote(tt.args.cu, tt.args.note)
		})
	}
}

func Test_createFormRegistrationUser(t *testing.T) {
	type args struct {
		cu *UIController
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateFormRegistrationUser",
			args: args{
				cu: &UIController{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormRegistrationUser(tt.args.cu)
		})
	}
}

func Test_createFormTextNote(t *testing.T) {
	type args struct {
		cu   *UIController
		note models.TextNote
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test_createFormTextNote",
			args{
				cu:   &UIController{},
				note: models.TextNote{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createFormTextNote(tt.args.cu, tt.args.note)
		})
	}
}

func Test_createMainMenu(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCreateMainMenu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createMainMenu()
		})
	}
}

func Test_createModalError(t *testing.T) {
	type args struct {
		err          error
		switchToPage string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateModalError",
			args: args{
				err:          errors.New("TestCreateModalError"),
				switchToPage: PageMenu,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createModalError(tt.args.err, tt.args.switchToPage)
		})
	}
}

func Test_createModalForm(t *testing.T) {
	type args struct {
		p      tview.Primitive
		width  int
		height int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateModalForm",
			args: args{
				p:      tview.NewFlex(),
				width:  10,
				height: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createModalForm(tt.args.p, tt.args.width, tt.args.height)
			assert.NotNil(t, got)
		})
	}
}

func Test_createNotesList(t *testing.T) {
	type args struct {
		storage []models.Noteable
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCreateNotesList",
			args: args{
				storage: []models.Noteable{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createNotesList(tt.args.storage)
		})
	}
}

func Test_creteMainFlex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestCreateMainFlex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := creteMainFlex()
			assert.NotNil(t, got)
		})
	}
}

func Test_formatCardNumber(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "16 number characters",
			args: args{
				text: "1234123412341234",
			},
			want: "1234-1234-1234-1234",
		},
		{
			name: "15 number characters",
			args: args{
				text: "123412341234123",
			},
			want: "1234-1234-1234-123",
		},
		{
			name: "15 number characters with alphabetic characters",
			args: args{
				text: "D1F2AS341234SS123F41G2G3H",
			},
			want: "D1F2-AS34-1234-SS12-3F41-G2G3-H",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatCardNumber(tt.args.text); got != tt.want {
				t.Errorf("formatCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
