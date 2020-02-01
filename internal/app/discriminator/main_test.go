package discriminator

import "testing"

func Test_stringMapEquals(t *testing.T) {
	tests := []struct {
		name string
		a    map[string]string
		b    map[string]string
		want bool
	}{
		{
			name: "unset value not same as set",
			a: map[string]string{
				"1": "",
			},
			b: map[string]string{
				"2": "",
			},
			want: false,
		},
		{
			name: "same",
			a: map[string]string{
				"1": "",
			},
			b: map[string]string{
				"1": "",
			},
			want: true,
		},
		{
			name: "different length",
			a: map[string]string{
				"1": "",
			},
			b:    map[string]string{},
			want: false,
		},
		{
			name: "b nil",
			a: map[string]string{
				"1": "",
			},
			b:    nil,
			want: false,
		},
		{
			name: "a nil",
			a:    nil,
			b: map[string]string{
				"1": "",
			},
			want: false,
		},
		{
			name: "both nil",
			a:    nil,
			b:    nil,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringMapEquals(tt.a, tt.b); got != tt.want {
				t.Errorf("stringMapEquals() = %v, want %v", got, tt.want)
			}
		})
	}
}
