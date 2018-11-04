package usecase

import "testing"

func TestInitConfig_BuildSpec(t *testing.T) {
	cases := []struct {
		test string
		cfg  InitConfig
		out  string
	}{
		{
			test: "empty",
		},
		{
			test: "HEAD",
			cfg:  InitConfig{HEAD: true},
			out:  "@master",
		},
		{
			test: "branch",
			cfg:  InitConfig{Branch: "foo/bar"},
			out:  "@foo/bar",
		},
		{
			test: "version",
			cfg:  InitConfig{Version: "^0.3.0"},
			out:  "@^0.3.0",
		},
		{
			test: "revision",
			cfg:  InitConfig{Revision: "a2489d2"},
			out:  "@a2489d2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			if got, want := tc.cfg.BuildSpec(), tc.out; got != want {
				t.Errorf("BuildSpec() returned %q, want %q", got, want)
			}
		})
	}
}
