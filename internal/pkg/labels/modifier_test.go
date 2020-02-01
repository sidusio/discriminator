package labels

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

func TestModifier(t *testing.T) {

	tests := []struct {
		name         string
		inputLabels  map[string]string
		outputLabels map[string]string
		modifier     []byte
	}{
		{
			name:        "multiple equal signs",
			inputLabels: map[string]string{},
			outputLabels: map[string]string{
				"test": "123=abc",
			},
			modifier: []byte("+test=123=abc"),
		},
		{
			name: "remove and add same label",
			inputLabels: map[string]string{
				"test": "123",
			},
			outputLabels: map[string]string{},
			modifier:     []byte("+test=abc\n-test"),
		},
		{
			name: "padded -",
			inputLabels: map[string]string{
				"test": "123",
			},
			outputLabels: map[string]string{},
			modifier:     []byte("\n     -test \n"),
		},
		{
			name:        "padded +",
			inputLabels: map[string]string{},
			outputLabels: map[string]string{
				"test": "123",
			},
			modifier: []byte("\n     +test=123 \n"),
		},
		{
			name:        "messy",
			inputLabels: map[string]string{},
			outputLabels: map[string]string{
				"test": "123",
			},
			modifier: []byte("sdfsdfg\n     +test=123 \n sdfsdf"),
		},
		{
			name: "update label",
			inputLabels: map[string]string{
				"test": "123",
			},
			outputLabels: map[string]string{
				"test": "abc",
			},
			modifier: []byte("+test=abc"),
		},
		{
			name: "remove label",
			inputLabels: map[string]string{
				"test": "123",
			},
			outputLabels: map[string]string{},
			modifier:     []byte("-test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewModifier(context.Background(), bytes.NewReader(tt.modifier))
			if err != nil {
				t.Errorf("NewModifier() error = %v", err)
				return
			}
			got := m.Apply(tt.inputLabels)
			if !reflect.DeepEqual(got, tt.outputLabels) {
				t.Errorf("Labels got = %v, want %v", got, tt.outputLabels)
			}
		})
	}
}

func TestModifiers_Apply(t *testing.T) {
	tests := []struct {
		name   string
		ms     Modifiers
		labels map[string]string
		want   map[string]string
	}{
		{
			name: "order check",
			ms: Modifiers{
				{
					deletions: []string{"a"},
				},
				{
					additions: map[string]string{"a": "1"},
				},
			},
			labels: map[string]string{},
			want:   map[string]string{"a": "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ms.Apply(tt.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
