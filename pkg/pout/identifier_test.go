package pout

import "testing"

func TestSplitIdentifier(t *testing.T) {
	t.Parallel()
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "no sub-package",
			args: args{
				input: "test_pkg.Message",
			},
			want:    "test_pkg",
			want1:   "Message",
			wantErr: false,
		},
		{
			name: "single sub-package",
			args: args{
				input: "test_pkg.test_sub_pkg.Message",
			},
			want:    "test_pkg.test_sub_pkg",
			want1:   "Message",
			wantErr: false,
		},
		{
			name: "google timestamp package",
			args: args{
				input: "google.protobuf.Timestamp",
			},
			want:    "google.protobuf",
			want1:   "Timestamp",
			wantErr: false,
		},
		{
			name: "does not allow full stops at the end",
			args: args{
				input: "test_pkg.Message.",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := SplitIdentifier(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SplitIdentifier() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitIdentifier() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
