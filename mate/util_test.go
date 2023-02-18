package mate

import (
	"fmt"
	"github.com/ericnts/log"
	"reflect"
	"testing"
)

type A interface {
}

type B struct {
	a A
}

func TestA(t *testing.T) {
	aType := reflect.TypeOf((*A)(nil))
	aName := fmt.Sprintf("%s.%s", aType.PkgPath(), aType.Name())
	log.Info(aName)
	log.Info()

}

func TestParseMethodName(t *testing.T) {
	type args struct {
		methodName string
	}
	tests := []struct {
		name       string
		args       args
		wantMethod string
		wantPath   string
	}{
		{
			name: "1",
			args: args{
				methodName: "GetSprintsItemInfo",
			},
			wantMethod: "GET",
			wantPath:   "/sprints/:id/info",
		},
		{
			name: "2",
			args: args{
				methodName: "GetPlansItem",
			},
			wantMethod: "GET",
			wantPath:   "/plans/:id",
		},
		{
			name: "3",
			args: args{
				methodName: "GetItemQRCode",
			},
			wantMethod: "GET",
			wantPath:   "/:id/qrCode",
		},
		{
			name: "4",
			args: args{
				methodName: "PutEnablesOfRooms",
			},
			wantMethod: "PUT",
			wantPath:   "/rooms/enables",
		},
		{
			name: "5",
			args: args{
				methodName: "PutEnablesOfRoomsItem",
			},
			wantMethod: "PUT",
			wantPath:   "/rooms/:id/enables",
		},
		{
			name: "6",
			args: args{
				methodName: "PutRoomEnablesOfItem",
			},
			wantMethod: "PUT",
			wantPath:   "/:id/roomEnables",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotPath := parseMethodName(tt.args.methodName)
			if gotMethod != tt.wantMethod {
				t.Errorf("parseMethodName() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseMethodName() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func Test_parseMethodName(t *testing.T) {
	type args struct {
		methodName string
	}
	tests := []struct {
		name       string
		args       args
		wantMethod string
		wantPath   string
	}{
		{
			args: args{
				methodName: "PostPicturesItemOfItem",
			},
			wantMethod: "POST",
			wantPath:   "/:id/pictures/:id2",
		},
		{
			args: args{
				methodName: "GetItemCategories",
			},
			wantMethod: "GET",
			wantPath:   "/:id/categories",
		},
		{
			args: args{
				methodName: "GetSerialsItemCategories",
			},
			wantMethod: "GET",
			wantPath:   "/serials/:id/categories",
		},
		{
			args: args{
				methodName: "PostItemPictures",
			},
			wantMethod: "POST",
			wantPath:   "/:id/pictures",
		},
		{
			args: args{
				methodName: "PostItemPicturesItem",
			},
			wantMethod: "POST",
			wantPath:   "/:id/pictures/:id2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotPath := parseMethodName(tt.args.methodName)
			if gotMethod != tt.wantMethod {
				t.Errorf("parseMethodName() gotMethod = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Errorf("parseMethodName() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
