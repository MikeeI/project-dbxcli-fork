package cmd

import (
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/spf13/cobra"
)

const (
	searchFlagNameGlob                   = "name-glob"
	searchFlagNameRegex                  = "name-regex"
	searchFlagNameContains               = "name-contains"
	searchFlagExcludeNameContains        = "exclude-name-contains"
	searchFlagMinSize                    = "min-size"
	searchFlagMaxSize                    = "max-size"
	searchFlagModifiedFrom               = "modified-from"
	searchFlagModifiedBefore             = "modified-before"
	searchFlagType                       = "type"
	searchEntryTypeAny                   = ""
	searchEntryTypeFile                  = "file"
	searchEntryTypeFolder                = "folder"
	searchDefaultFilterText              = ""
	searchDefaultFilterBytes      uint64 = 0
)

var searchEntryTypes = []string{searchEntryTypeFile, searchEntryTypeFolder}

type searchFilters struct {
	nameGlob            string
	nameRegex           *regexp.Regexp
	nameContains        string
	excludeNameContains []string
	minSize             *uint64
	maxSize             *uint64
	modifiedFrom        *time.Time
	modifiedBefore      *time.Time
	entryType           string
	timeField           string
}

func registerSearchFilterFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.String(searchFlagNameGlob, searchDefaultFilterText, "Local case-sensitive glob matching the entire entry name (not its path)")
	flags.String(searchFlagNameRegex, searchDefaultFilterText, "Local Go regular expression on the entry name; use (?i) for case-insensitive matching")
	flags.String(searchFlagNameContains, searchDefaultFilterText, "Local case-sensitive substring required in the entry name")
	flags.StringArray(searchFlagExcludeNameContains, nil, "Exclude entry names containing this case-sensitive substring (repeatable; commas are literal)")
	flags.Uint64(searchFlagMinSize, searchDefaultFilterBytes, "Local inclusive minimum file size in bytes (unset means no lower bound)")
	flags.Uint64(searchFlagMaxSize, searchDefaultFilterBytes, "Local inclusive maximum file size in bytes (unset means no upper bound; 0 means empty files)")
	flags.String(searchFlagModifiedFrom, searchDefaultFilterText, "Local inclusive file modification lower bound in RFC3339; uses --time")
	flags.String(searchFlagModifiedBefore, searchDefaultFilterText, "Local exclusive file modification upper bound in RFC3339; uses --time")
	flags.String(searchFlagType, searchEntryTypeAny, "Local entry type filter: file, folder (unset allows both)")
}

func parseSearchFilters(cmd *cobra.Command, timeField string) (searchFilters, error) {
	filters := searchFilters{timeField: timeField}
	flags := cmd.Flags()
	var pattern string
	for _, field := range []struct {
		name   string
		target *string
	}{
		{searchFlagNameGlob, &filters.nameGlob},
		{searchFlagNameRegex, &pattern},
		{searchFlagNameContains, &filters.nameContains},
		{searchFlagType, &filters.entryType},
	} {
		value, err := flags.GetString(field.name)
		if err != nil {
			return searchFilters{}, err
		}
		*field.target = value
	}
	if _, err := path.Match(filters.nameGlob, ""); err != nil {
		return searchFilters{}, invalidArgumentsErrorfWithDetails("invalid --%s %q: %v", flagValueErrorDetails(searchFlagNameGlob, filters.nameGlob), searchFlagNameGlob, filters.nameGlob, err)
	}
	if pattern != "" {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return searchFilters{}, invalidArgumentsErrorfWithDetails("invalid --%s %q: %v", flagValueErrorDetails(searchFlagNameRegex, pattern), searchFlagNameRegex, pattern, err)
		}
		filters.nameRegex = compiled
	}
	excluded, err := flags.GetStringArray(searchFlagExcludeNameContains)
	if err != nil {
		return searchFilters{}, err
	}
	for _, value := range excluded {
		if value == "" {
			return searchFilters{}, invalidArgumentsErrorWithDetails("--exclude-name-contains must not be empty", flagValueErrorDetails(searchFlagExcludeNameContains, value))
		}
	}
	filters.excludeNameContains = excluded
	for _, field := range []struct {
		name   string
		target **uint64
	}{
		{searchFlagMinSize, &filters.minSize},
		{searchFlagMaxSize, &filters.maxSize},
	} {
		if !flags.Changed(field.name) {
			continue
		}
		value, err := flags.GetUint64(field.name)
		if err != nil {
			return searchFilters{}, err
		}
		*field.target = &value
	}
	for _, field := range []struct {
		name   string
		target **time.Time
	}{
		{searchFlagModifiedFrom, &filters.modifiedFrom},
		{searchFlagModifiedBefore, &filters.modifiedBefore},
	} {
		if !flags.Changed(field.name) {
			continue
		}
		value, err := flags.GetString(field.name)
		if err != nil {
			return searchFilters{}, err
		}
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return searchFilters{}, invalidArgumentsErrorfWithDetails("invalid --%s %q: use RFC3339 with a timezone", flagValueErrorDetails(field.name, value), field.name, value)
		}
		*field.target = &parsed
	}
	if filters.entryType != searchEntryTypeAny && filters.entryType != searchEntryTypeFile && filters.entryType != searchEntryTypeFolder {
		return searchFilters{}, invalidArgumentsErrorfWithDetails("invalid --type %q: use file or folder", flagValueErrorDetails(searchFlagType, filters.entryType), filters.entryType)
	}
	if filters.minSize != nil && filters.maxSize != nil && *filters.minSize > *filters.maxSize {
		return searchFilters{}, invalidArgumentsErrorWithDetails("--min-size must not exceed --max-size", flagsErrorDetails(searchFlagMinSize, searchFlagMaxSize))
	}
	if filters.modifiedFrom != nil && filters.modifiedBefore != nil && !filters.modifiedFrom.Before(*filters.modifiedBefore) {
		return searchFilters{}, invalidArgumentsErrorWithDetails("--modified-from must precede --modified-before", flagsErrorDetails(searchFlagModifiedFrom, searchFlagModifiedBefore))
	}
	if filters.entryType == searchEntryTypeFolder && filters.hasFileBounds() {
		return searchFilters{}, invalidArgumentsErrorWithDetails("--type folder cannot be combined with file size or modification bounds", flagsErrorDetails(searchFlagType, searchFlagMinSize, searchFlagMaxSize, searchFlagModifiedFrom, searchFlagModifiedBefore))
	}
	return filters, nil
}

func (f searchFilters) hasFileBounds() bool {
	return f.minSize != nil || f.maxSize != nil || f.modifiedFrom != nil || f.modifiedBefore != nil
}

func (f searchFilters) active() bool {
	return f.nameGlob != "" || f.nameRegex != nil || f.nameContains != "" || len(f.excludeNameContains) > 0 || f.entryType != searchEntryTypeAny || f.hasFileBounds()
}

func (f searchFilters) matches(metadata files.IsMetadata) (bool, error) {
	name := metadataName(metadata)
	if f.nameGlob != "" {
		matched, err := path.Match(f.nameGlob, name)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	if f.nameRegex != nil && !f.nameRegex.MatchString(name) {
		return false, nil
	}
	if !strings.Contains(name, f.nameContains) {
		return false, nil
	}
	for _, excluded := range f.excludeNameContains {
		if strings.Contains(name, excluded) {
			return false, nil
		}
	}
	file, isFile := metadata.(*files.FileMetadata)
	if f.entryType == searchEntryTypeFile && !isFile || f.entryType == searchEntryTypeFolder && isFile {
		return false, nil
	}
	if !f.hasFileBounds() {
		return true, nil
	}
	if !isFile {
		return false, nil
	}
	if f.minSize != nil && file.Size < *f.minSize || f.maxSize != nil && file.Size > *f.maxSize {
		return false, nil
	}
	if f.modifiedFrom != nil || f.modifiedBefore != nil {
		modified := getTime(file, listOptions{timeField: f.timeField})
		if modified.IsZero() {
			return false, commandFailedErrorfWithDetails("search metadata for %q has no %s modification time", operationErrorDetails("search"), file.Id, f.timeField)
		}
		if f.modifiedFrom != nil && modified.Before(*f.modifiedFrom) || f.modifiedBefore != nil && !modified.Before(*f.modifiedBefore) {
			return false, nil
		}
	}
	return true, nil
}
