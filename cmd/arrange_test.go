package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/n3dst4/gopotato/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	oldJournalDate = "2020-01-15"
	oldMonthFolder = "2020-01"
)

func setupJournals(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	journals := filepath.Join(root, "journals")
	require.NoError(t, os.MkdirAll(journals, 0755))
	config = &utils.Config{
		RootPath:     root,
		JournalsPath: journals,
		PagesPath:    filepath.Join(root, "pages"),
		KeepDays:     7,
		KeepMonths:   6,
	}
	t.Chdir(root)
	return journals
}

func writeJournal(t *testing.T, journals, date, content string) string {
	t.Helper()
	path := filepath.Join(journals, date+".md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestArrange_CreatesScratchFile(t *testing.T) {
	journals := setupJournals(t)
	arrange()
	assert.FileExists(t, filepath.Join(journals, "scratch.md"))
}

func TestArrange_ScratchFileIsIdempotent(t *testing.T) {
	journals := setupJournals(t)
	scratchPath := filepath.Join(journals, "scratch.md")
	require.NoError(t, os.WriteFile(scratchPath, []byte("my scratch notes"), 0644))
	arrange()
	content, err := os.ReadFile(scratchPath)
	require.NoError(t, err)
	assert.Equal(t, "my scratch notes", string(content))
}

func TestArrange_CreatesJournalFile(t *testing.T) {
	journals := setupJournals(t)
	arrange()
	today := time.Now().Format("2006-01-02")
	journalPath := filepath.Join(journals, today+".md")
	assert.FileExists(t, journalPath)
	content, err := os.ReadFile(journalPath)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("# %s\n\n", today), string(content))
}

func TestArrange_JournalFileIsIdempotent(t *testing.T) {
	journals := setupJournals(t)
	today := time.Now().Format("2006-01-02")
	existingContent := fmt.Sprintf("# %s\n\nsome actual content", today)
	writeJournal(t, journals, today, existingContent)
	arrange()
	content, err := os.ReadFile(filepath.Join(journals, today+".md"))
	require.NoError(t, err)
	assert.Equal(t, existingContent, string(content))
}

func TestArrange_ArchivesOldJournalWithContent(t *testing.T) {
	journals := setupJournals(t)
	original := writeJournal(t, journals, oldJournalDate, fmt.Sprintf("# %s\n\nsome real content", oldJournalDate))
	arrange()
	assert.NoFileExists(t, original)
	// archiveOldJournals moves it to YYYY-MM/, then archiveOldMonths moves that to YYYY/YYYY-MM/
	assert.FileExists(t, filepath.Join(journals, "2020", oldMonthFolder, oldJournalDate+".md"))
}

func TestArrange_DeletesUnmodifiedOldJournal(t *testing.T) {
	journals := setupJournals(t)
	original := writeJournal(t, journals, oldJournalDate, fmt.Sprintf("# %s\n\n", oldJournalDate))
	arrange()
	assert.NoFileExists(t, original)
	assert.NoDirExists(t, filepath.Join(journals, oldMonthFolder))
}

func TestArrange_LeavesRecentJournalAlone(t *testing.T) {
	journals := setupJournals(t)
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	recentPath := writeJournal(t, journals, yesterday, fmt.Sprintf("# %s\n\nsome content", yesterday))
	arrange()
	assert.FileExists(t, recentPath)
}

func TestArrange_ScratchNotArchivedOrDeleted(t *testing.T) {
	journals := setupJournals(t)
	scratchPath := filepath.Join(journals, "scratch.md")
	require.NoError(t, os.WriteFile(scratchPath, []byte("scratch content"), 0644))
	arrange()
	content, err := os.ReadFile(scratchPath)
	require.NoError(t, err)
	assert.Equal(t, "scratch content", string(content))
}

func TestArrange_ArchivesOldMonthFolder(t *testing.T) {
	journals := setupJournals(t)
	oldMonthDir := filepath.Join(journals, oldMonthFolder)
	require.NoError(t, os.MkdirAll(oldMonthDir, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(oldMonthDir, oldJournalDate+".md"),
		[]byte(fmt.Sprintf("# %s\n\nsome content", oldJournalDate)),
		0644,
	))
	arrange()
	assert.NoDirExists(t, oldMonthDir)
	assert.FileExists(t, filepath.Join(journals, "2020", oldMonthFolder, oldJournalDate+".md"))
}
