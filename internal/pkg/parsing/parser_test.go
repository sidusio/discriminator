package parsing

import (
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		template  string
		arguments map[string]string
		wantErr   bool
	}{
		{
			name:      "no arguments",
			input:     "a()",
			template:  "a",
			arguments: map[string]string{},
		},
		{
			name:    "invalid",
			input:   "a(:)",
			wantErr: true,
		},
		{
			name:    "nothing",
			input:   "",
			wantErr: true,
		},
		{
			name:     "one argument",
			input:    "a(b:4)",
			template: "a",
			arguments: map[string]string{
				"b": "4",
			},
		},
		{
			name:     "two arguments",
			input:    "a(b:4, c: 7)",
			template: "a",
			arguments: map[string]string{
				"b": "4",
				"c": "7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, arguments, err := parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if template != tt.template {
				t.Errorf("parse() template = %v, want %v", template, tt.template)
			}
			if !reflect.DeepEqual(arguments, tt.arguments) {
				t.Errorf("parse() arguments = %v, want %v", arguments, tt.arguments)
			}
		})
	}
}
