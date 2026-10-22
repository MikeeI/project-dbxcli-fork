// Copyright © 2016 Dropbox, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dropbox/dbxcli/v3/internal/jsonschema"
	"github.com/dropbox/dbxcli/v3/internal/output"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/spf13/cobra"
)

type searchInput struct {
	Query               string     `json:"query"`
	Path                string     `json:"path,omitempty"`
	Content             bool       `json:"content"`
	Limit               uint64     `json:"limit,omitempty"`
	OrderBy             string     `json:"order_by,omitempty"`
	Long                bool       `json:"long"`
	Sort                string     `json:"sort,omitempty"`
	Reverse             bool       `json:"reverse"`
	Time                string     `json:"time,omitempty"`
	TimeFormat          string     `json:"time_format,omitempty"`
	Category            []string   `json:"category,omitempty"`
	Ext                 []string   `json:"ext,omitempty"`
	Highlights          bool       `json:"highlights,omitempty"`
	NameGlob            string     `json:"name_glob,omitempty"`
	NameRegex           string     `json:"name_regex,omitempty"`
	NameContains        string     `json:"name_contains,omitempty"`
	ExcludeNameContains []string   `json:"exclude_name_contains,omitempty"`
	MinSize             *uint64    `json:"min_size,omitempty"`
	MaxSize             *uint64    `json:"max_size,omitempty"`
	ModifiedFrom        *time.Time `json:"modified_from,omitempty"`
	ModifiedBefore      *time.Time `json:"modified_before,omitempty"`
	Type                string     `json:"type,omitempty"`
}

const (
	searchJSONStatusFound                = "found"
	searchStopAPIExhausted               = "api_exhausted"
	searchStopLimit                      = "limit"
	searchStopError                      = "error"
	searchSummaryDetailKey               = "search"
	searchPartialResultsDetailKey        = "partial_results"
	searchFlagCategory                   = "category"
	searchFlagExt                        = "ext"
	searchFlagHighlights                 = "highlights"
	searchDefaultContent                 = false
	searchDefaultLimit            uint64 = 0
	searchDefaultOrderBy                 = ""
	searchDefaultLong                    = false
	searchDefaultSort                    = ""
	searchDefaultReverse                 = false
	searchDefaultTime                    = "server"
	searchDefaultTimeFormat              = ""
	searchDefaultHighlights              = false
)

var searchCategories = []string{
	files.FileCategoryImage, files.FileCategoryDocument, files.FileCategoryPdf,
	files.FileCategorySpreadsheet, files.FileCategoryPresentation, files.FileCategoryAudio,
	files.FileCategoryVideo, files.FileCategoryFolder, files.FileCategoryPaper, files.FileCategoryOthers,
}

type searchCommandOptions struct {
	list       listOptions
	content    bool
	limit      uint64
	orderBy    string
	categories []*files.FileCategory
	extensions []string
	highlights bool
	filters    searchFilters
}

type searchResults struct {
	matches        []*files.SearchMatchV2
	serverHasMore  *bool
	stopReason     string
	apiMatchesSeen uint64
}

type searchSummary struct {
	Limit               uint64 `json:"limit"`
	StopReason          string `json:"stop_reason"`
	ServerHasMore       *bool  `json:"server_has_more"`
	LocalFiltersApplied bool   `json:"local_filters_applied"`
	APIMatchesSeen      uint64 `json:"api_matches_seen"`
	APIMatchLimit       uint64 `json:"api_match_limit"`
}

type jsonSearchMetadata struct {
	jsonMetadata
	Name           string                 `json:"name,omitempty"`
	MatchType      string                 `json:"match_type,omitempty"`
	HighlightSpans []*files.HighlightSpan `json:"highlight_spans,omitempty"`
}

type jsonSearchOutput struct {
	jsonOperationOutput
	Search searchSummary `json:"search"`
}

func search(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return invalidArgumentsErrorWithDetails("`search` requires a `query` argument", argumentErrorDetails("query"))
	}
	if len(args) > 2 {
		return invalidArgumentsErrorWithDetails("`search` accepts at most one optional `path-scope` argument", mergeJSONErrorDetails(argumentErrorDetails("path-scope"), pathErrorDetails(args[2])))
	}

	var scope string
	if len(args) == 2 {
		scope = args[1]
		if !strings.HasPrefix(scope, "/") {
			return invalidArgumentsErrorWithDetails("`search` `path-scope` must begin with \"/\"", mergeJSONErrorDetails(argumentErrorDetails("path-scope"), pathErrorDetails(scope)))
		}
	}

	opts, err := parseSearchOptions(cmd)
	if err != nil {
		return err
	}
	arg := newSearchV2Arg(args[0], scope, opts)

	results, searchErr := collectSearchMatches(currentContext(), filesNewFunc(config), arg, opts.filters, opts.limit)
	if opts.list.sortBy != "" {
		sort.SliceStable(results.matches, func(i, j int) bool {
			return compareLess(results.matches[i].Metadata.Metadata, results.matches[j].Metadata.Metadata, opts.list)
		})
	}
	return renderSearchOutput(cmd, args[0], scope, results, opts, searchErr)
}

func parseSearchOptions(cmd *cobra.Command) (searchCommandOptions, error) {
	content, err := cmd.Flags().GetBool("content")
	if err != nil {
		return searchCommandOptions{}, err
	}
	limit, err := cmd.Flags().GetUint64("limit")
	if err != nil {
		return searchCommandOptions{}, err
	}
	orderBy, err := cmd.Flags().GetString("order-by")
	if err != nil {
		return searchCommandOptions{}, err
	}
	listOpts, err := parseListOptions(cmd)
	if err != nil {
		return searchCommandOptions{}, err
	}
	if !validSearchOrderBy(orderBy) {
		return searchCommandOptions{}, invalidArgumentsErrorWithDetails("`search --order-by` must be one of: relevance, modified", flagErrorDetails("order-by"))
	}
	highlights, err := cmd.Flags().GetBool(searchFlagHighlights)
	if err != nil {
		return searchCommandOptions{}, err
	}
	if highlights && commandOutputFormat(cmd) != output.FormatJSON {
		return searchCommandOptions{}, invalidArgumentsErrorWithDetails("--highlights requires --output=json", flagsErrorDetails(searchFlagHighlights, outputFlag))
	}
	categories, err := cmd.Flags().GetStringSlice(searchFlagCategory)
	if err != nil {
		return searchCommandOptions{}, err
	}
	var fileCategories []*files.FileCategory
	for _, category := range categories {
		if !slices.Contains(searchCategories, category) {
			return searchCommandOptions{}, invalidArgumentsErrorfWithDetails("invalid --category %q", flagValueErrorDetails(searchFlagCategory, category), category)
		}
		fileCategories = append(fileCategories, &files.FileCategory{Tagged: dropbox.Tagged{Tag: category}})
	}
	extensions, err := cmd.Flags().GetStringSlice(searchFlagExt)
	if err != nil {
		return searchCommandOptions{}, err
	}
	for i, extension := range extensions {
		normalized := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(extension), "."))
		if normalized == "" || strings.ContainsAny(normalized, "*?[]/\\") {
			return searchCommandOptions{}, invalidArgumentsErrorfWithDetails("invalid --ext %q: use extensions without glob patterns", flagValueErrorDetails(searchFlagExt, extension), extension)
		}
		extensions[i] = normalized
	}
	filters, err := parseSearchFilters(cmd, listOpts.timeField)
	if err != nil {
		return searchCommandOptions{}, err
	}

	return searchCommandOptions{
		list:       listOpts,
		content:    content,
		limit:      limit,
		orderBy:    orderBy,
		categories: fileCategories,
		extensions: extensions,
		highlights: highlights,
		filters:    filters,
	}, nil
}

func newSearchV2Arg(query, scope string, opts searchCommandOptions) *files.SearchV2Arg {
	arg := files.NewSearchV2Arg(query)
	searchOpts := files.NewSearchOptions()
	searchOpts.Path = scope
	searchOpts.FilenameOnly = !opts.content
	if !opts.filters.active() && opts.limit > 0 && opts.limit < searchOpts.MaxResults {
		searchOpts.MaxResults = opts.limit
	}
	if opts.orderBy != "" {
		searchOpts.OrderBy = searchOrderBy(opts.orderBy)
	}
	searchOpts.FileCategories = opts.categories
	searchOpts.FileExtensions = opts.extensions
	if opts.highlights {
		arg.MatchFieldOptions = &files.SearchMatchFieldOptions{IncludeHighlights: true}
	}
	arg.Options = searchOpts
	return arg
}

func validSearchOrderBy(orderBy string) bool {
	switch orderBy {
	case "", "relevance", "modified":
		return true
	default:
		return false
	}
}

func searchOrderBy(orderBy string) *files.SearchOrderBy {
	switch orderBy {
	case "modified":
		return &files.SearchOrderBy{Tagged: dropbox.Tagged{Tag: files.SearchOrderByLastModifiedTime}}
	default:
		return &files.SearchOrderBy{Tagged: dropbox.Tagged{Tag: files.SearchOrderByRelevance}}
	}
}

// collectSearchMatches returns every accepted match even when a later page fails.
// Only accepted IDs consume the limit; rejected versions must not hide a later qualifying version.
func collectSearchMatches(ctx context.Context, client filesClient, arg *files.SearchV2Arg, filters searchFilters, limit uint64) (searchResults, error) {
	results := searchResults{stopReason: searchStopError}
	seen := make(map[string]struct{})
	page, err := client.SearchV2Context(ctx, arg)
	for {
		if err != nil {
			return results, err
		}
		if page == nil {
			return results, errors.New("search returned a nil result page")
		}
		results.serverHasMore = new(page.HasMore)
		results.apiMatchesSeen += uint64(len(page.Matches))
		omitted := false
		for _, match := range page.Matches {
			metadata, id, err := searchMatchMetadata(match)
			if err != nil {
				return results, err
			}
			accepted, err := filters.matches(metadata)
			if err != nil {
				return results, err
			}
			if !accepted {
				continue
			}
			if _, duplicate := seen[id]; duplicate {
				continue
			}
			if limit > 0 && uint64(len(results.matches)) >= limit {
				omitted = true
				continue
			}
			seen[id] = struct{}{}
			results.matches = append(results.matches, match)
		}
		// Inspect the rest of this page: an exact final-page limit is not truncation.
		// HasMore alone only proves unvisited pages, not additional locally matching IDs.
		if omitted || page.HasMore && limit > 0 && uint64(len(results.matches)) >= limit {
			results.stopReason = searchStopLimit
			return results, nil
		}
		if !page.HasMore {
			results.stopReason = searchStopAPIExhausted
			return results, nil
		}
		if page.Cursor == "" {
			return results, errors.New("search returned has_more without a continuation cursor")
		}
		page, err = client.SearchContinueV2Context(ctx, files.NewSearchV2ContinueArg(page.Cursor))
	}
}

func searchMatchMetadata(match *files.SearchMatchV2) (files.IsMetadata, string, error) {
	if match == nil || match.Metadata == nil || match.Metadata.Tag != files.MetadataV2Metadata {
		return nil, "", errors.New("search returned unsupported or missing match metadata")
	}
	var id string
	switch metadata := match.Metadata.Metadata.(type) {
	case *files.FileMetadata:
		if metadata != nil {
			id = metadata.Id
		}
	case *files.FolderMetadata:
		if metadata != nil {
			id = metadata.Id
		}
	default:
		return nil, "", fmt.Errorf("search returned unsupported metadata type %T", match.Metadata.Metadata)
	}
	if id == "" {
		return nil, "", errors.New("search returned file or folder metadata without an ID")
	}
	for _, span := range match.HighlightSpans {
		if span == nil {
			return nil, "", fmt.Errorf("search returned a nil title highlight span for %q", id)
		}
	}
	return match.Metadata.Metadata, id, nil
}

func newSearchSummary(results searchResults, opts searchCommandOptions) searchSummary {
	return searchSummary{
		Limit:               opts.limit,
		StopReason:          results.stopReason,
		ServerHasMore:       results.serverHasMore,
		LocalFiltersApplied: opts.filters.active(),
		APIMatchesSeen:      results.apiMatchesSeen,
		APIMatchLimit:       jsonschema.SearchAPIMatchLimit,
	}
}

func searchOperationResults(matches []*files.SearchMatchV2) ([]jsonOperationResult, error) {
	results := make([]jsonOperationResult, 0, len(matches))
	for _, match := range matches {
		metadata, err := jsonMetadataFromDropbox(match.Metadata.Metadata)
		if err != nil {
			return nil, err
		}
		result := jsonSearchMetadata{
			jsonMetadata:   metadata,
			Name:           metadataName(match.Metadata.Metadata),
			HighlightSpans: match.HighlightSpans,
		}
		if match.MatchType != nil {
			result.MatchType = match.MatchType.Tag
		}
		results = append(results, newJSONOperationResult(searchJSONStatusFound, metadata.Type, nil, result))
	}
	return results, nil
}

func renderSearchOutput(cmd *cobra.Command, query, scope string, results searchResults, opts searchCommandOptions, searchErr error) error {
	summary := newSearchSummary(results, opts)
	out := commandOutput(cmd)
	if commandOutputFormat(cmd) != output.FormatJSON {
		entries := make([]files.IsMetadata, 0, len(results.matches))
		for _, match := range results.matches {
			entries = append(entries, match.Metadata.Metadata)
		}
		writeErr := out.RenderText(func(w io.Writer) error { return renderSearchResults(w, entries, opts.list) })
		if writeErr == nil {
			serverHasMore := "unknown"
			if summary.ServerHasMore != nil {
				serverHasMore = strconv.FormatBool(*summary.ServerHasMore)
			}
			_, writeErr = fmt.Fprintf(cmd.ErrOrStderr(), "search stop_reason=%s limit=%d returned=%d server_has_more=%s local_filters_applied=%t api_matches_seen=%d api_match_limit=%d\n",
				summary.StopReason, summary.Limit, len(entries), serverHasMore, summary.LocalFiltersApplied, summary.APIMatchesSeen, summary.APIMatchLimit)
		}
		if writeErr != nil && searchErr != nil {
			// Joining alone hides SDK summaries from the legacy single-cause extractor.
			return newCodedError(jsonErrorCode(searchErr), errors.Join(searchErr, writeErr), jsonErrorDetails(searchErr))
		}
		if writeErr != nil {
			return writeErr
		}
		return searchErr
	}
	operationResults, err := searchOperationResults(results.matches)
	if err != nil {
		return err
	}
	if searchErr != nil {
		// Root owns the sole error document; do not emit a success envelope first.
		return withJSONErrorDetails(searchErr, map[string]any{
			searchSummaryDetailKey:        summary,
			searchPartialResultsDetailKey: operationResults,
		})
	}
	payload := jsonSearchOutput{
		jsonOperationOutput: newJSONCommandOperationOutput(cmd, newSearchInput(query, scope, opts), operationResults, nil),
		Search:              summary,
	}
	return out.Render(nil, payload)
}

func newSearchInput(query, scope string, opts searchCommandOptions) searchInput {
	input := searchInput{
		Query:               query,
		Path:                scope,
		Content:             opts.content,
		Limit:               opts.limit,
		OrderBy:             opts.orderBy,
		Long:                opts.list.long,
		Sort:                opts.list.sortBy,
		Reverse:             opts.list.reverse,
		Time:                opts.list.timeField,
		TimeFormat:          opts.list.timeFormat,
		Ext:                 opts.extensions,
		Highlights:          opts.highlights,
		NameGlob:            opts.filters.nameGlob,
		NameContains:        opts.filters.nameContains,
		ExcludeNameContains: opts.filters.excludeNameContains,
		MinSize:             opts.filters.minSize,
		MaxSize:             opts.filters.maxSize,
		ModifiedFrom:        opts.filters.modifiedFrom,
		ModifiedBefore:      opts.filters.modifiedBefore,
		Type:                opts.filters.entryType,
	}
	for _, category := range opts.categories {
		input.Category = append(input.Category, category.Tag)
	}
	if opts.filters.nameRegex != nil {
		input.NameRegex = opts.filters.nameRegex.String()
	}
	return input
}

func renderSearchResults(out io.Writer, entries []files.IsMetadata, opts listOptions) error {
	w := new(tabwriter.Writer)
	w.Init(out, 4, 8, 1, ' ', 0)

	if opts.long {
		if _, err := fmt.Fprint(w, "Revision\tSize\tLast modified\tPath\n"); err != nil {
			return err
		}
	}

	for _, entry := range entries {
		switch f := entry.(type) {
		case *files.FileMetadata:
			if _, err := fmt.Fprintln(w, formatFileMetadataWithOpts(f, opts)); err != nil {
				return err
			}
		case *files.FolderMetadata:
			var row strings.Builder
			printFolderMetadata(&row, f, opts.long)
			if _, err := fmt.Fprintln(w, row.String()); err != nil {
				return err
			}
		}
	}

	return w.Flush()
}

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [flags] <query> [path-scope]",
	Short: "Search",
	RunE:  search,
}

func init() {
	RootCmd.AddCommand(searchCmd)
	registerSearchFlags(searchCmd)
	enableStructuredOutput(searchCmd)
}

func registerSearchFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.BoolP("content", "c", searchDefaultContent, "Search file contents in addition to filenames")
	flags.Uint64("limit", searchDefaultLimit, "Maximum distinct matching entries to return; applied before local sorting (0 is unlimited)")
	flags.String("order-by", searchDefaultOrderBy, "Server-side search ordering: relevance, modified")
	flags.BoolP("long", "l", searchDefaultLong, "Long listing")
	flags.String("sort", searchDefaultSort, "Sort the selected results locally by: name, size, time, type; not a global top-N selection")
	flags.BoolP("reverse", "r", searchDefaultReverse, "Reverse the chosen local sort order")
	flags.String("time", searchDefaultTime, "Time field for display, local sorting and modification filters: server, client")
	flags.String("time-format", searchDefaultTimeFormat, "Time format: short (2006-01-02 15:04), rfc3339")
	flags.StringSlice(searchFlagCategory, nil, "Server-side active-file categories (comma-separated)")
	flags.StringSlice(searchFlagExt, nil, "Server-side active-file extensions, e.g. pdf,docx (comma-separated; no glob patterns)")
	flags.Bool(searchFlagHighlights, searchDefaultHighlights, "Include title highlight spans; requires --output=json (not content snippets)")
	registerSearchFilterFlags(cmd)
	for name, values := range map[string][]string{
		searchFlagCategory: searchCategories,
		searchFlagType:     searchEntryTypes,
	} {
		// The command's fixed registration must not hide duplicate-completion programmer errors.
		if err := cmd.RegisterFlagCompletionFunc(name, cobra.FixedCompletions(values, cobra.ShellCompDirectiveNoFileComp)); err != nil {
			panic(err)
		}
	}
}
