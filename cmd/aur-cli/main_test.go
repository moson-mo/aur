package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/Jguer/aur"
	"github.com/stretchr/testify/assert"
)

func Test_printSearch(t *testing.T) {
	var b bytes.Buffer

	a := &aur.Pkg{
		Name:        "test",
		Version:     "1.0.0.",
		NumVotes:    20,
		Popularity:  4.0,
		Description: "Test description",
	}

	testWriter := io.Writer(&b)

	printSearch(a, testWriter)

	assert.Equal(t, "- \x1b[1mtest\x1b[0m 1.0.0. (20 4.00)\n\tTest description\n", b.String())
}

func Test_printInfo(t *testing.T) {
	a := &aur.Pkg{
		Name:        "test",
		Version:     "1.0.0.",
		NumVotes:    20,
		Popularity:  4.0,
		Description: "Test description",
	}

	tests := []struct {
		name    string
		verbose bool
		wantW   string
	}{
		{
			name:    "verbose",
			verbose: true,
			wantW:   "\x1b[1mName            : \x1b[0mtest\n\x1b[1mVersion         : \x1b[0m1.0.0.\n\x1b[1mDescription     : \x1b[0mTest description\n\x1b[1mKeywords        : \x1b[0mNone\n\x1b[1mURL             : \x1b[0mNone\n\x1b[1mAUR URL         : \x1b[0mhttps://aur.archlinux.org/packages/test\n\x1b[1mGroups          : \x1b[0mNone\n\x1b[1mLicenses        : \x1b[0mNone\n\x1b[1mProvides        : \x1b[0mNone\n\x1b[1mDepends On      : \x1b[0mNone\n\x1b[1mMake Deps       : \x1b[0mNone\n\x1b[1mCheck Deps      : \x1b[0mNone\n\x1b[1mOptional Deps   : \x1b[0mNone\n\x1b[1mConflicts With  : \x1b[0mNone\n\x1b[1mMaintainer      : \x1b[0mNone\n\x1b[1mVotes           : \x1b[0m20\n\x1b[1mPopularity      : \x1b[0m4.000000\n\x1b[1mFirst Submitted : \x1b[0mThu 01 Jan 1970 01:00:00 AM CET\n\x1b[1mLast Modified   : \x1b[0mThu 01 Jan 1970 01:00:00 AM CET\n\x1b[1mOut-of-date     : \x1b[0mNo\n\x1b[1mID              : \x1b[0m0\n\x1b[1mPackage Base ID : \x1b[0m0\n\x1b[1mPackage Base    : \x1b[0mNone\n\x1b[1mSnapshot URL    : \x1b[0mhttps://aur.archlinux.org\n\n",
		},
		{
			name:    "not verbose",
			verbose: false,
			wantW:   "\x1b[1mName            : \x1b[0mtest\n\x1b[1mVersion         : \x1b[0m1.0.0.\n\x1b[1mDescription     : \x1b[0mTest description\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			testWriter := io.Writer(&b)

			printInfo(a, testWriter, "https://aur.archlinux.org", tt.verbose)
			assert.Equal(t, tt.wantW, b.String())
		})
	}

}
