package utils

import (
	"reflect"
	"testing"
)

func TestWildCardToRegexp(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				pattern: "test",
			},
			want: "^" + "(?i)" + "test" + "$",
		},
		{
			name: "test*",
			args: args{
				pattern: "test*",
			},
			want: "^" + "(?i)" + "test.*" + "$",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildCardToRegexp(tt.args.pattern); got != tt.want {
				t.Errorf("WildCardToRegexp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWildCardMatch(t *testing.T) {
	type args struct {
		pattern string
		value   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{
				pattern: "test",
				value:   "test",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildCardMatch(tt.args.pattern, tt.args.value); got != tt.want {
				t.Errorf("WildCardMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWildCardMatchs(t *testing.T) {
	type args struct {
		pattern []string
		value   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{
				pattern: []string{"test"},
				value:   "test",
			},
			want: true,
		},
		{
			name: "test*",
			args: args{
				pattern: []string{"test*"},
				value:   "test",
			},
			want: true,
		},
		{
			name: "test*",
			args: args{
				pattern: []string{"test*"},
				value:   "test1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WildCardMatchs(tt.args.pattern, tt.args.value); got != tt.want {
				t.Errorf("WildCardMatchs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchPathFiles(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				filename: "test",
			},
			want:    []string{"/usr/bin/test", "/bin/test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchPathFiles(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPathFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchPathFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
