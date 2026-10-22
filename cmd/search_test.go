package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/spf13/cobra"
)

func TestSearchArgValidation(t *testing.T) {
	err := search(searchCmd, []string{})
	if err == nil {
		t.Error("expected error for no args")
	}
}

func TestSearchPathScopeValidation(t *testing.T) {
	err := search(searchCmd, []string{"query", "no-slash"})
	if err == nil {
		t.Fatal("expected error for path-scope without leading slash")
	}
	details := jsonErrorDetails(err)
	if details["argument"] != "path-scope" || details["path"] != "no-slash" {
		t.Fatalf("details = %#v, want path-scope argument and path", details)
	}
}

func TestSearchRejectsExtraArgs(t *testing.T) {
	err := search(searchCmd, []string{"query", "/docs", "extra"})
	if err == nil || !strings.Contains(err.Error(), "path-scope") {
		t.Fatalf("error = %v, want extra path-scope error", err)
	}
	details := jsonErrorDetails(err)
	if details["argument"] != "path-scope" || details["path"] != "extra" {
		t.Fatalf("details = %#v, want extra path-scope argument path", details)
	}
}

func TestSearchOrderByValidation(t *testing.T) {
	cmd, _ := testSearchCmd()
	setSearchFlag(t, cmd, "order-by", "name")

	err := search(cmd, []string{"query"})
	if err == nil || !strings.Contains(err.Error(), "order-by") {
		t.Fatalf("error = %v, want order-by validation error", err)
	}
}

func TestSearchRejectsInvalidListOptions(t *testing.T) {
	tests := []struct {
		name  string
		flag  string
		value string
	}{
		{name: "sort", flag: "sort", value: "date"},
		{name: "time", flag: "time", value: "created"},
		{name: "time-format", flag: "time-format", value: "unix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, stdout := testSearchCmd()
			setSearchFlag(t, cmd, tt.flag, tt.value)
			stubFilesClient(t, &mockFilesClient{
				searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
					t.Fatalf("SearchV2 called for invalid --%s", tt.flag)
					return nil, nil
				},
			})

			err := search(cmd, []string{"query"})
			if err == nil {
				t.Fatalf("expected invalid --%s error", tt.flag)
			}
			if !strings.Contains(err.Error(), tt.flag) || !strings.Contains(err.Error(), tt.value) {
				t.Fatalf("error = %v, want flag and value", err)
			}
			if got := stdout.String(); got != "" {
				t.Fatalf("stdout = %q, want empty output on validation error", got)
			}
		})
	}
}

func TestRenderSearchResultsSeparatesMatchesWithNewlines(t *testing.T) {
	entries := []files.IsMetadata{
		&files.FileMetadata{
			Metadata: files.Metadata{PathDisplay: "/first.txt"},
		},
		&files.FolderMetadata{
			Metadata: files.Metadata{PathDisplay: "/second"},
		},
	}

	var out bytes.Buffer
	if err := renderSearchResults(&out, entries, listOptions{long: false}); err != nil {
		t.Fatalf("renderSearchResults returned error: %v", err)
	}

	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected one output line per match, got %d lines in %q", len(lines), out.String())
	}

	if got, want := strings.TrimSpace(lines[0]), "/first.txt"; got != want {
		t.Errorf("first rendered match = %q, want %q", got, want)
	}
	if got, want := strings.TrimSpace(lines[1]), "/second"; got != want {
		t.Errorf("second rendered match = %q, want %q", got, want)
	}
}

func TestRenderSearchResultsLongModeIncludesHeader(t *testing.T) {
	entries := []files.IsMetadata{
		&files.FileMetadata{
			Metadata: files.Metadata{PathDisplay: "/first.txt"},
			Rev:      "abc123",
			Size:     42,
		},
	}

	var out bytes.Buffer
	if err := renderSearchResults(&out, entries, listOptions{long: true}); err != nil {
		t.Fatalf("renderSearchResults returned error: %v", err)
	}

	got := out.String()
	for _, want := range []string{"Revision", "Size", "Last modified", "Path", "abc123", "/first.txt"} {
		if !strings.Contains(got, want) {
			t.Errorf("output = %q, want to contain %q", got, want)
		}
	}
}

func TestSearchUsesSearchV2AndCommandOutput(t *testing.T) {
	cmd, stdout := testSearchCmd()
	for name, value := range map[string]string{
		"limit": "2", searchFlagNameGlob: "*.pdf", searchFlagNameRegex: "(?i)report",
		searchFlagNameContains: "Report", searchFlagExcludeNameContains: "Draft,old",
		searchFlagMinSize: "0", searchFlagMaxSize: "42", searchFlagType: searchEntryTypeFile,
		searchFlagModifiedFrom: "2026-01-01T00:00:00Z", searchFlagModifiedBefore: "2027-01-01T00:00:00Z",
		"time": "client",
	} {
		setSearchFlag(t, cmd, name, value)
	}
	first := &files.FileMetadata{
		Metadata: files.Metadata{Name: "Report-first.pdf", PathDisplay: "/docs/Report-first.pdf"},
		Id:       "id:first", ClientModified: dropbox.DBXTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	later := &files.FileMetadata{
		Metadata: files.Metadata{Name: "Report-second.pdf", PathDisplay: "/docs/Report-second.pdf"},
		Id:       "id:version", Size: 42, ClientModified: first.ClientModified,
	}
	rejected := *later
	rejected.Name = "Report-Draft,old.pdf"
	var firstArg *files.SearchV2Arg
	var continueCursor string
	stubFilesClient(t, &mockFilesClient{
		searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
			firstArg = arg
			page := files.NewSearchV2Result([]*files.SearchMatchV2{
				searchMatch(&rejected),
				searchMatch(&files.FolderMetadata{Metadata: files.Metadata{Name: "Report-folder.pdf"}, Id: "id:folder"}),
				searchMatch(first),
			}, true)
			page.Cursor = "cursor-1"
			return page, nil
		},
		searchContinueV2Fn: func(arg *files.SearchV2ContinueArg) (*files.SearchV2Result, error) {
			continueCursor = arg.Cursor
			return files.NewSearchV2Result([]*files.SearchMatchV2{searchMatch(first), searchMatch(later)}, false), nil
		},
	})
	if err := search(cmd, []string{"needle", "/docs"}); err != nil {
		t.Fatalf("search: %v", err)
	}
	if firstArg == nil || firstArg.Options == nil || firstArg.Query != "needle" || firstArg.Options.Path != "/docs" || !firstArg.Options.FilenameOnly {
		t.Fatalf("first request = %#v", firstArg)
	}
	if firstArg.Options.MaxResults != files.NewSearchOptions().MaxResults {
		t.Fatalf("local filters reduced API page size to %d", firstArg.Options.MaxResults)
	}
	if continueCursor != "cursor-1" {
		t.Fatalf("continuation cursor = %q", continueCursor)
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 || strings.TrimSpace(lines[0]) != first.PathDisplay || strings.TrimSpace(lines[1]) != later.PathDisplay {
		t.Fatalf("filtered distinct results = %q", stdout.String())
	}
	status := cmd.ErrOrStderr().(*bytes.Buffer).String()
	for _, want := range []string{"stop_reason=api_exhausted", "returned=2", "server_has_more=false", "local_filters_applied=true", "api_matches_seen=5"} {
		if !strings.Contains(status, want) {
			t.Fatalf("status = %q, missing %q", status, want)
		}
	}
}

func TestSearchJSONOutputsInputAndResults(t *testing.T) {
	cmd, stdout := testSearchCmd()
	setSearchOutputJSON(t, cmd)
	setSearchFlag(t, cmd, "long", "true")
	setSearchFlag(t, cmd, "sort", "name")
	setSearchFlag(t, cmd, "reverse", "true")
	setSearchFlag(t, cmd, "time", "client")
	setSearchFlag(t, cmd, "time-format", "rfc3339")
	setSearchFlag(t, cmd, "content", "true")
	setSearchFlag(t, cmd, "limit", "2")
	setSearchFlag(t, cmd, "order-by", "modified")
	setSearchFlag(t, cmd, searchFlagCategory, "document,folder")
	setSearchFlag(t, cmd, searchFlagExt, ".TXT,docx")
	setSearchFlag(t, cmd, searchFlagHighlights, "true")
	var firstArg *files.SearchV2Arg

	mock := &mockFilesClient{
		searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
			firstArg = arg
			matches := []*files.SearchMatchV2{
				searchMatch(&files.FileMetadata{
					Metadata: files.Metadata{
						Name:        "report.txt",
						PathDisplay: "/docs/report.txt",
						PathLower:   "/docs/report.txt",
					},
					Id:   "id:file",
					Rev:  "rev-file",
					Size: 42,
				}),
				searchMatch(&files.FolderMetadata{
					Metadata: files.Metadata{
						Name:        "archive",
						PathDisplay: "/docs/archive",
						PathLower:   "/docs/archive",
					},
					Id: "id:folder",
				}),
			}
			matches[0].MatchType = &files.SearchMatchTypeV2{Tagged: dropbox.Tagged{Tag: files.SearchMatchTypeV2FileContent}}
			matches[0].HighlightSpans = []*files.HighlightSpan{files.NewHighlightSpan("report", true), files.NewHighlightSpan(".txt", false)}
			return files.NewSearchV2Result(matches, false), nil
		},
	}
	stubFilesClient(t, mock)

	if err := search(cmd, []string{"report", "/docs"}); err != nil {
		t.Fatalf("search error: %v", err)
	}
	if firstArg == nil {
		t.Fatal("SearchV2 was not called")
	}
	if firstArg.Query != "report" {
		t.Fatalf("query = %q, want report", firstArg.Query)
	}
	if firstArg.Options == nil || firstArg.Options.Path != "/docs" {
		t.Fatalf("options path = %#v, want /docs", firstArg.Options)
	}
	if firstArg.Options.FilenameOnly {
		t.Fatal("filename_only = true, want false with --content")
	}
	if firstArg.Options.MaxResults != 2 {
		t.Fatalf("max_results = %d, want 2", firstArg.Options.MaxResults)
	}
	if firstArg.Options.OrderBy == nil || firstArg.Options.OrderBy.Tag != files.SearchOrderByLastModifiedTime {
		t.Fatalf("order_by = %#v, want last_modified_time", firstArg.Options.OrderBy)
	}
	if len(firstArg.Options.FileCategories) != 2 || firstArg.Options.FileCategories[0].Tag != files.FileCategoryDocument ||
		len(firstArg.Options.FileExtensions) != 2 || firstArg.Options.FileExtensions[0] != "txt" ||
		firstArg.MatchFieldOptions == nil || !firstArg.MatchFieldOptions.IncludeHighlights || firstArg.IncludeHighlights {
		t.Fatalf("server filter/highlight options = %#v", firstArg)
	}

	got := decodeSearchOutput(t, stdout)
	if got.Input.Query != "report" || got.Input.Path != "/docs" {
		t.Fatalf("input = %#v, want query report path /docs", got.Input)
	}
	if !got.Input.Content || got.Input.Limit != 2 || got.Input.OrderBy != "modified" || !got.Input.Long || got.Input.Sort != "name" || !got.Input.Reverse || got.Input.Time != "client" || got.Input.TimeFormat != "rfc3339" {
		t.Fatalf("input options = %#v, want content/limit/order-by/long/sort/reverse/time/time-format", got.Input)
	}
	if len(got.Results) != 2 {
		t.Fatalf("results = %d, want 2", len(got.Results))
	}
	first := got.Results[0].Result
	if got.Results[0].Status != searchJSONStatusFound || got.Results[0].Kind != "file" || first.Type != "file" || first.Rev != "rev-file" || first.Size == nil || *first.Size != 42 {
		t.Fatalf("first result = %#v, want found file metadata", got.Results[0])
	}
	if first.Name != "report.txt" || first.MatchType != files.SearchMatchTypeV2FileContent ||
		len(first.HighlightSpans) != 2 || !first.HighlightSpans[0].IsHighlighted {
		t.Fatalf("search match evidence = %#v", first)
	}
	if got.Search.StopReason != searchStopAPIExhausted || got.Search.Limit != 2 ||
		got.Search.ServerHasMore == nil || *got.Search.ServerHasMore || got.Search.APIMatchesSeen != 2 {
		t.Fatalf("exact final-page limit = %#v", got.Search)
	}
	if !got.Input.Highlights || len(got.Input.Category) != 2 || len(got.Input.Ext) != 2 || got.Input.Ext[0] != "txt" {
		t.Fatalf("search filter input = %#v", got.Input)
	}
	second := got.Results[1].Result
	if got.Results[1].Status != searchJSONStatusFound || got.Results[1].Kind != "folder" || second.Type != "folder" || second.ID != "id:folder" {
		t.Fatalf("second result = %#v, want found folder metadata", got.Results[1])
	}
	if strings.Contains(stdout.String(), `"entries"`) {
		t.Fatalf("JSON output = %s, want operation results and no entries key", stdout.String())
	}
}

func TestSearchJSONOmitsPathWithoutScope(t *testing.T) {
	cmd, stdout := testSearchCmd()
	setSearchOutputJSON(t, cmd)

	mock := &mockFilesClient{
		searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
			if arg.Options == nil {
				t.Fatal("options = nil, want search options")
			}
			if arg.Options.Path != "" {
				t.Fatalf("options path = %q, want empty", arg.Options.Path)
			}
			if !arg.Options.FilenameOnly {
				t.Fatal("filename_only = false, want true by default")
			}
			return files.NewSearchV2Result(nil, false), nil
		},
	}
	stubFilesClient(t, mock)

	if err := search(cmd, []string{"report"}); err != nil {
		t.Fatalf("search error: %v", err)
	}
	output := append([]byte(nil), stdout.Bytes()...)
	got := decodeSearchOutput(t, stdout)
	if got.Input.Query != "report" || got.Input.Path != "" {
		t.Fatalf("input = %#v, want query report and empty path", got.Input)
	}
	var raw map[string]any
	if err := json.Unmarshal(output, &raw); err != nil {
		t.Fatalf("decode raw JSON output: %v\noutput: %s", err, string(output))
	}
	input, ok := raw["input"].(map[string]any)
	if !ok {
		t.Fatalf("raw input = %#v, want object", raw["input"])
	}
	if _, ok := input["path"]; ok {
		t.Fatalf("input path key is present in %s, want omitted", string(output))
	}
}

func TestSearchLimitCapsResultsAndStopsPagination(t *testing.T) {
	for _, test := range []struct {
		name     string
		extra    bool
		hasMore  bool
		wantStop string
	}{
		{name: "unvisited pages", hasMore: true, wantStop: searchStopLimit},
		{name: "unreturned last-page match", extra: true, wantStop: searchStopLimit},
		{name: "exact final page with duplicate", wantStop: searchStopAPIExhausted},
	} {
		t.Run(test.name, func(t *testing.T) {
			cmd, stdout := testSearchCmd()
			setSearchOutputJSON(t, cmd)
			setSearchFlag(t, cmd, "limit", "1")
			first := searchMatch(&files.FileMetadata{Metadata: files.Metadata{Name: "first.txt", PathDisplay: "/docs/first.txt"}, Id: "id:first"})
			stubFilesClient(t, &mockFilesClient{
				searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
					if arg.Options.MaxResults != 1 {
						t.Fatalf("unfiltered page size = %d", arg.Options.MaxResults)
					}
					matches := []*files.SearchMatchV2{first, first}
					if test.extra {
						matches = append(matches, searchMatch(&files.FileMetadata{Metadata: files.Metadata{Name: "second.txt"}, Id: "id:second"}))
					}
					page := files.NewSearchV2Result(matches, test.hasMore)
					page.Cursor = "next"
					return page, nil
				},
				searchContinueV2Fn: func(arg *files.SearchV2ContinueArg) (*files.SearchV2Result, error) {
					t.Fatalf("unexpected lookahead after limit: %s", arg.Cursor)
					return nil, nil
				},
			})
			if err := search(cmd, []string{"needle", "/docs"}); err != nil {
				t.Fatalf("search: %v", err)
			}
			got := decodeSearchOutput(t, stdout)
			if len(got.Results) != 1 || got.Results[0].Result.ID != "id:first" ||
				got.Search.StopReason != test.wantStop || got.Search.ServerHasMore == nil || *got.Search.ServerHasMore != test.hasMore {
				t.Fatalf("bounded search = %#v", got)
			}
		})
	}
}

func TestSearchJSONErrorWritesNoOutput(t *testing.T) {
	for _, continuation := range []bool{false, true} {
		t.Run(fmt.Sprintf("continuation=%t", continuation), func(t *testing.T) {
			cmd, stdout := testSearchCmd()
			setSearchOutputJSON(t, cmd)
			cause := files.SearchContinueV2APIError{APIError: dropbox.APIError{ErrorSummary: "path/not_found/"}}
			mock := &mockFilesClient{
				searchV2Fn: func(arg *files.SearchV2Arg) (*files.SearchV2Result, error) {
					if !continuation {
						return nil, cause
					}
					page := files.NewSearchV2Result([]*files.SearchMatchV2{searchMatch(&files.FileMetadata{
						Metadata: files.Metadata{Name: "report.pdf", PathDisplay: "/report.pdf"},
						Id:       "id:report",
					})}, true)
					page.Cursor = "next"
					return page, nil
				},
				searchContinueV2Fn: func(arg *files.SearchV2ContinueArg) (*files.SearchV2Result, error) {
					return nil, cause
				},
			}
			stubFilesClient(t, mock)
			err := search(cmd, []string{"report"})
			if err == nil || jsonErrorCode(err) != jsonErrorCodeNotFound || exitCodeForError(err) != exitCodeNotFound {
				t.Fatalf("error = %v, want original not-found code and exit status", err)
			}
			if got := stdout.String(); got != "" {
				t.Fatalf("stdout = %q, want no document before root error rendering", got)
			}
			renderCommandError(cmd, err)
			var response struct {
				OK    bool `json:"ok"`
				Error struct {
					Code    string `json:"code"`
					Details struct {
						Search struct {
							StopReason    string `json:"stop_reason"`
							ServerHasMore *bool  `json:"server_has_more"`
						} `json:"search"`
						PartialResults []jsonOperationResult `json:"partial_results"`
					} `json:"details"`
				} `json:"error"`
			}
			// Unmarshal rejects concatenated documents, protecting the root-renderer boundary.
			if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
				t.Fatalf("decode single error document: %v\n%s", err, stdout.String())
			}
			wantResults := 0
			if continuation {
				wantResults = 1
			}
			if response.OK || response.Error.Code != jsonErrorCodeNotFound ||
				response.Error.Details.Search.StopReason != "error" ||
				len(response.Error.Details.PartialResults) != wantResults {
				t.Fatalf("partial error response = %#v", response)
			}
			hasMore := response.Error.Details.Search.ServerHasMore
			if continuation && (hasMore == nil || !*hasMore) || !continuation && hasMore != nil {
				t.Fatalf("last successful server state = %v, continuation=%t", hasMore, continuation)
			}
		})
	}
}

func TestSearchCommandSupportsStructuredOutput(t *testing.T) {
	if !commandSupportsStructuredOutput(searchCmd) {
		t.Fatal("search command should support structured output")
	}
}

func testSearchCmd() (*cobra.Command, *bytes.Buffer) {
	var stdout bytes.Buffer
	cmd := &cobra.Command{Use: "search"}
	cmd.SetOut(&stdout)
	cmd.SetErr(&bytes.Buffer{})
	registerSearchFlags(cmd)
	cmd.Flags().String(outputFlag, "text", "")
	return cmd, &stdout
}

func setSearchOutputJSON(t *testing.T, cmd *cobra.Command) {
	t.Helper()
	setSearchFlag(t, cmd, outputFlag, "json")
}

func setSearchFlag(t *testing.T, cmd *cobra.Command, name, value string) {
	t.Helper()
	if err := cmd.Flags().Set(name, value); err != nil {
		t.Fatalf("set %s: %v", name, err)
	}
}

type searchOutputForTest struct {
	Input   searchInput   `json:"input"`
	Search  searchSummary `json:"search"`
	Results []struct {
		Status string             `json:"status"`
		Kind   string             `json:"kind"`
		Result jsonSearchMetadata `json:"result"`
	} `json:"results"`
}

func decodeSearchOutput(t *testing.T, out *bytes.Buffer) searchOutputForTest {
	t.Helper()
	var result searchOutputForTest
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("decode search: %v\n%s", err, out.String())
	}
	return result
}

func searchMatch(metadata files.IsMetadata) *files.SearchMatchV2 {
	return files.NewSearchMatchV2(&files.MetadataV2{
		Tagged:   dropbox.Tagged{Tag: files.MetadataV2Metadata},
		Metadata: metadata,
	})
}
