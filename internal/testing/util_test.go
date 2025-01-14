package testing_test

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	internaltesting "github.com/squeedee/geo/internal/testing"
	"testing"
)

var oneRepeatOutput = `words
and
pictures
and
sounds`

type TestMock struct {
	FatalfCallCount   int
	FatalfCallStrings []string
	HelperCallCount   int
}

func (m *TestMock) Fatalf(format string, args ...any) {
	m.FatalfCallCount += 1
	m.FatalfCallStrings = append(m.FatalfCallStrings, fmt.Sprintf(format, args...))
}

func (m *TestMock) Helper() {
	m.HelperCallCount += 1
}

func TestStringMatcher_MatchText(t *testing.T) {
	tests := map[string]struct {
		output         string
		matchText      []string
		shouldFailWith string
	}{
		"matches progressively": {
			output:    oneRepeatOutput,
			matchText: []string{"words", "and", "pict", "ures", "and", "sounds"},
		},
		"no failure for one matched instance": {
			output:    oneRepeatOutput,
			matchText: []string{"and"},
		},
		"no failure for two matched instance": {
			output:    oneRepeatOutput,
			matchText: []string{"and", "and"},
		},
		"failure for last unmatched instance": {
			output:         oneRepeatOutput,
			matchText:      []string{"and", "and", "and"},
			shouldFailWith: "could not match 'and' in remaining output:\n\nsounds\n",
		},
		"failure for first query passing the second": {
			output:         oneRepeatOutput,
			matchText:      []string{"and", "words"},
			shouldFailWith: "could not match 'words' in remaining output:\n\npictures\nand\nsounds\n",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mock := &TestMock{}
			matcher := internaltesting.NewOutputMatcher(tc.output)

			for _, match := range tc.matchText {
				matcher.MatchText(mock, match)
			}

			if mock.HelperCallCount != len(tc.matchText) {
				t.Fatalf("Expected helper to be called %d time(s), call count: %d", len(tc.matchText), mock.HelperCallCount)
			}

			if mock.FatalfCallCount == 0 {
				if tc.shouldFailWith != "" {
					t.Fatalf("did not fail with expected: %s ", tc.shouldFailWith)
				}
			} else {
				if tc.shouldFailWith == "" {
					t.Fatalf("did not expect failure, got %d failures, first one: %s ",
						mock.FatalfCallCount,
						mock.FatalfCallStrings[0])
				} else if diff := cmp.Diff(tc.shouldFailWith, mock.FatalfCallStrings[0]); diff != "" {
					t.Fatalf("should have failed with diff: \n%s", diff)
				}
			}
		})
	}
}
