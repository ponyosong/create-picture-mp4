package main

import "testing"

func TestDisposeSpecialCharacter(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "testCase01", args: args{s: `Don'\''t Stay - Linkin Park.mp3`}, want: `Don't Stay - Linkin Park.mp3`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DisposeSpecialCharacter(tt.args.s); got != tt.want {
				t.Errorf("DisposeSpecialCharacter() = %v, want %v", got, tt.want)
			}
		})
	}
}
