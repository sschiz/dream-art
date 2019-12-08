package shop

import (
	"reflect"
	"testing"

	"github.com/sschiz/dream-art/pkg/category"
)

func TestShop_AppendAdmin(t *testing.T) {
	type fields struct {
		categories []category.Category
		admins     []string
		Syncer     *Syncer
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "WrongNickname",
			fields:  fields{admins: []string{"lol"}},
			args:    args{"e"},
			want:    []string{"lol"},
			wantErr: true,
		},
		{
			name:    "AlreadyExists",
			fields:  fields{admins: []string{"lol"}},
			args:    args{"@lol"},
			want:    []string{"lol"},
			wantErr: true,
		},
		{
			name:    "CommonplaceCase",
			fields:  fields{admins: []string{"a", "c"}},
			args:    args{"@b"},
			want:    []string{"a", "b", "c"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shop{
				categories: tt.fields.categories,
				admins:     tt.fields.admins,
				Syncer:     tt.fields.Syncer,
			}
			if err := s.AppendAdmin(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Shop.AppendAdmin() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(s.admins, tt.want) {
				t.Errorf("got = %v, expected = %v", s.admins, tt.want)
			}
		})
	}
}
