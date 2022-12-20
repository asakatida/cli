package lock

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/cli/v2/internal/prompter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/httpmock"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/cli/v2/test"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func Test_NewCmdLock(t *testing.T) {
	// TODO parent name stuff?
	cases := []struct {
		name    string
		args    string
		want    LockOptions
		wantErr string
		tty     bool
	}{
		{
			name: "sets reason",
			args: "--reason off_topic 451",
			want: LockOptions{
				Reason:      "off_topic",
				SelectorArg: "451",
			},
		},
		{
			name:    "no args",
			wantErr: "accepts 1 arg(s), received 0",
		},
		{
			name: "no flags",
			args: "451",
			want: LockOptions{
				SelectorArg: "451",
			},
		},
		{
			name:    "bad reason",
			args:    "--reason bad 451",
			wantErr: "invalid reason bad",
		},
		{
			name:    "bad reason tty",
			args:    "--reason bad 451",
			tty:     true,
			wantErr: "X Invalid reason: bad\n",
		},
		{
			name: "interactive",
			args: "451",
			tty:  true,
			want: LockOptions{
				SelectorArg: "451",
				Interactive: true,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			ios.SetStdoutTTY(tt.tty)
			ios.SetStdinTTY(tt.tty)
			ios.SetStderrTTY(tt.tty)
			f := &cmdutil.Factory{
				IOStreams: ios,
			}
			var opts *LockOptions
			cmd := NewCmdLock(f, "issue", func(_ string, o *LockOptions) error {
				opts = o
				return nil
			})
			cmd.PersistentFlags().StringP("repo", "R", "", "")

			argv, err := shlex.Split(tt.args)
			assert.NoError(t, err)

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.Reason, opts.Reason)
			assert.Equal(t, tt.want.SelectorArg, opts.SelectorArg)
			assert.Equal(t, tt.want.Interactive, opts.Interactive)
		})
	}
}

func Test_NewCmdUnlock(t *testing.T) {
	// TODO parent name stuff?
	cases := []struct {
		name    string
		args    string
		want    LockOptions
		wantErr string
		tty     bool
	}{
		{
			name:    "no args",
			wantErr: "accepts 1 arg(s), received 0",
		},
		{
			name: "no flags",
			args: "451",
			want: LockOptions{
				SelectorArg: "451",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()
			ios.SetStdoutTTY(tt.tty)
			ios.SetStdinTTY(tt.tty)
			ios.SetStderrTTY(tt.tty)
			f := &cmdutil.Factory{
				IOStreams: ios,
			}
			var opts *LockOptions
			cmd := NewCmdUnlock(f, "issue", func(_ string, o *LockOptions) error {
				opts = o
				return nil
			})
			cmd.PersistentFlags().StringP("repo", "R", "", "")

			argv, err := shlex.Split(tt.args)
			assert.NoError(t, err)

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
				return
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want.SelectorArg, opts.SelectorArg)
		})
	}
}

func Test_runLock(t *testing.T) {
	cases := []struct {
		name        string
		opts        LockOptions
		promptStubs func(*testing.T, *prompter.PrompterMock)
		httpStubs   func(*testing.T, *httpmock.Registry)
		wantOut     string
		wantErrOut  string
		wantErr     string
		tty         bool
		state       string
	}{
		{
			name: "lock an issue",
			// TODO
		},
		// TODO
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			reg := &httpmock.Registry{}
			reg.StubRepoInfoResponse("OWNER", "REPO", "trunk")
			defer reg.Verify(t)
			if tt.httpStubs != nil {
				tt.httpStubs(t, reg)
			}

			pm := &prompter.PrompterMock{}
			if tt.promptStubs != nil {
				tt.promptStubs(t, pm)
			}

			ios, _, stdout, stderr := iostreams.Test()
			// TODO do i need to bother with this
			ios.SetStdoutTTY(tt.tty)
			ios.SetStdinTTY(tt.tty)
			ios.SetStderrTTY(tt.tty)

			opts := LockOptions{
				Prompter: pm,
				IO:       ios,
				HttpClient: func() (*http.Client, error) {
					return &http.Client{Transport: reg}, nil
				},
			}

			err := lockRun(tt.state, &opts)
			output := &test.CmdOut{
				OutBuf: stdout,
				ErrBuf: stderr,
			}
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOut, output.String())
				assert.Equal(t, tt.wantErrOut, output.Stderr())
			}
		})
	}
}

func TestReasons(t *testing.T) {
	assert.Equal(t, len(reasons), len(reasonsApi))

	for _, reason := range reasons {
		assert.Equal(t, strings.ToUpper(reason), string(*reasonsMap[reason]))
	}
}
