package main

import "testing"

func TestParseMp3Duration(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "testCase-01", args: args{s: "00:00:30.01,"}, want: 30, wantErr: false},
		{name: "testCase-02", args: args{s: "00:00:30.01"}, want: 0, wantErr: true},
		{name: "testCase-03", args: args{s: "00:01:30.01,"}, want: 90, wantErr: false},
		{name: "testCase-04", args: args{s: "00:11:30.01,"}, want: 690, wantErr: false},
		{name: "testCase-05", args: args{s: "01:11:30.01,"}, want: 4290, wantErr: false},
		{name: "testCase-06", args: args{s: "00:00:00.01,"}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMp3Duration(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMp3Duration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMp3Duration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMp4Time(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{name: "testCase-01", args: args{s: "time=00:00:18.00"}, want: 18, wantErr: false},
		{name: "testCase-02", args: args{s: "00:00:30.01"}, want: 0, wantErr: true},
		{name: "testCase-03", args: args{s: "time=00:00:18.00 "}, want: 18, wantErr: false},
		{name: "testCase-04", args: args{s: "time=00:00:18.00    "}, want: 18, wantErr: false},
		{name: "testCase-05", args: args{s: "time=01:11:30.01,"}, want: 4290, wantErr: false},
		{name: "testCase-06", args: args{s: "time=00:00:00.00"}, want: 0, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMp4Time(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMp4Time() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseMp4Time() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecondsToStr(t *testing.T) {
	type args struct {
		t float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "testCase-01", args: args{t: 42.2}, want: "0小时0分钟42秒"},
		{name: "testCase-01", args: args{t: 4243.2}, want: "1小时10分钟43秒"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecondsToStr(tt.args.t); got != tt.want {
				t.Errorf("SecondsToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
