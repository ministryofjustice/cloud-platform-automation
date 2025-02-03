package utils

import (
	"testing"

	"github.com/google/go-github/v68/github"
)

func TestFiles(t *testing.T) {
	type args struct {
		sfname string
		files  []*github.CommitFile
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Pass",
			args: args{
				sfname: "APPLY_PIPELINE_SKIP_THIS_NAMESPACE",
				files: []*github.CommitFile{
					{
						Filename: github.String("APPLY_PIPELINE_SKIP_THIS_NAMESPACE"),
					},
				},
			},
			want: true,
		},
		{
			name: "Multiple File Fail",
			args: args{
				sfname: "APPLY_PIPELINE_SKIP_THIS_NAMESPACE",
				files: []*github.CommitFile{
					{
						Filename: github.String("APPLY_PIPELINE_SKIP_THIS_NAMESPACE"),
					},
					{
						Filename: github.String("README.md"),
					},
				},
			},
			want: false,
		},
		{
			name: "No File Fail",
			args: args{
				sfname: "APPLY_PIPELINE_SKIP_THIS_NAMESPACE",
				files: []*github.CommitFile{
					{
						Filename: github.String("README.md"),
					},
				},
			},
			want: false,
		},
		{
			name: "File extension fail",
			args: args{
				sfname: "APPLY_PIPELINE_SKIP_THIS_NAMESPACE.tf",
				files: []*github.CommitFile{
					{
						Filename: github.String("APPLY_PIPELINE_SKIP_THIS_NAMESPACE"),
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Files(tt.args.sfname, tt.args.files); got != tt.want {
				t.Errorf("Files() = %v, want %v", got, tt.want)
			}
		})
	}
}
