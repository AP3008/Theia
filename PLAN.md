# Non-Blocking Full Search for Large Directories

## Context
When entering full search mode (`/` twice), `CreateSystemFileList` recursively scans the entire directory tree **synchronously** inside `Update()`, freezing the TUI. The cached-keystroke fix prevents re-scanning on each keypress, but the **initial scan** still blocks. Additionally, `sortDirList` uses O(n²) bubble sort, compounding the problem on large result sets.

---

## Step 1: Replace bubble sort with `sort.Slice` — `filesystem/scanner.go`

Replace the `sortDirList` function (lines 89-103) with Go's standard library sort:

```go
import "sort" // add to imports

func sortDirList(dl []SystemFile) {
    sort.Slice(dl, func(i, j int) bool {
        return strings.ToLower(dl[i].Name) < strings.ToLower(dl[j].Name)
    })
}
```

This turns O(n²) sorting into O(n log n). Same comparison logic (case-insensitive).

---

## Step 2: Skip common large directories during recursion — `filesystem/scanner.go`

Add a package-level exclusion map:

```go
var excludedDirs = map[string]bool{
    "node_modules": true,
    "vendor":       true,
    "__pycache__":  true,
    "dist":         true,
    "target":       true,
    "build":        true,
    "venv":         true,
    ".venv":        true,
    ".cache":       true,
}
```

Modify the recursive guard at line 24 from:

```go
if value.IsDir() && !strings.HasPrefix(value.Name(), ".") {
```

to:

```go
if value.IsDir() && !strings.HasPrefix(value.Name(), ".") && !excludedDirs[value.Name()] {
```

---

## Step 3: Add `Loading` field to Model — `tui/model.go`

Add one field to the Model struct:

```go
Loading bool // true while async full scan is running
```

No change to `InitialModel` needed — zero-value `false` is correct.

---

## Step 4: Async scan with `tea.Cmd` — `tui/update.go`

### 4a. Define a custom message type (top of file):

```go
type fullScanResultMsg struct {
    files []filesystem.SystemFile
    err   error
}
```

### 4b. Define a `tea.Cmd` factory function:

```go
func runFullScanCmd(path string, showHidden bool) tea.Cmd {
    return func() tea.Msg {
        files, err := filesystem.CreateSystemFileList(path, showHidden, false, false, true)
        return fullScanResultMsg{files: files, err: err}
    }
}
```

Bubbletea runs the returned closure in a goroutine and delivers the `fullScanResultMsg` back to `Update()` when it completes.

### 4c. Replace the synchronous `case "/"` (full search entry):

Old:
```go
case "/":
    m.FullSearch = true
    cachedList, err := filesystem.CreateSystemFileList(...)
    ...
    m.FullSearchCache = cachedList
    m.SystemFiles = cachedList
```

New:
```go
case "/":
    m.FullSearch = true
    m.Loading = true
    m.SystemFiles = nil
    m.FullSearchCache = nil
    m.SearchInput.Reset()
    return m, runFullScanCmd(m.Path, m.Settings.ShowHidden)
```

The key difference: return a `tea.Cmd` instead of `nil`. This makes the scan non-blocking.

### 4d. Handle `fullScanResultMsg` in the top-level type switch:

Add this as a new case alongside the existing `case tea.KeyMsg:`:

```go
case fullScanResultMsg:
    m.Loading = false
    if msg.err != nil {
        m.FullSearch = false
        m.Searching = false
        originalList, err := filesystem.CreateSystemFileList(
            m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode, false,
        )
        if err != nil {
            return m, tea.Quit
        }
        m.SystemFiles = originalList
        return m, nil
    }
    // If user cancelled (pressed esc) while scan was running, discard result
    if !m.FullSearch {
        return m, nil
    }
    m.FullSearchCache = msg.files
    m.SystemFiles = msg.files
    m.Cursor = 0
    m.TopRow = 0
    return m, nil
```

### 4e. Add a loading guard at the start of `case tea.KeyMsg:`:

```go
case tea.KeyMsg:
    if m.Loading {
        switch msg.String() {
        case "esc", "ctrl+c":
            m.Loading = false
            m.FullSearch = false
            m.Searching = false
            m.SearchInput.Blur()
            m.SearchInput.Reset()
            originalList, err := filesystem.CreateSystemFileList(
                m.Path, m.Settings.ShowHidden, m.Settings.FileMode, m.Settings.DirMode, false,
            )
            if err != nil {
                return m, tea.Quit
            }
            m.SystemFiles = originalList
            m.Cursor = 0
            m.TopRow = 0
        }
        return m, nil
    }

    // ... existing three-branch key handling below ...
```

While loading: only esc/ctrl+c work (to cancel). All other keys are swallowed.

---

## Step 5: Loading indicator in View — `tui/view.go`

### In `normalView()`, before the `m.FullSearch` check:

```go
if m.Loading {
    loadingMsg := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#f6c177")).
        Render("  Scanning directory tree...")
    s.WriteString(loadingMsg + "\n")
    info := infoStyle.Render("\n [esc] cancel")
    s.WriteString(info)
    return s.String()
}
```

### In `longView()`, after the header (before the column headers):

Add the same loading check — render "Scanning directory tree..." and return early.

---

## Edge Cases to Handle

| Scenario | What happens |
|----------|-------------|
| Esc during loading | Restores normal listing. When goroutine result arrives later, `!m.FullSearch` guard discards it |
| Typing during loading | Keys are swallowed — nothing to filter yet |
| Spamming `/` | Loading guard prevents launching multiple goroutines |
| Empty scan result | Renders search bar with no results (existing behavior) |
| Scan error (permissions) | Exits full search, restores normal listing |

---

## Verification

1. `go build ./...` — must compile cleanly
2. Run on a small directory — normal nav, search, full search all work
3. Run on `$HOME` or a large project — `/ /` shows "Scanning..." immediately, UI stays responsive, Esc cancels
4. Confirm `node_modules`/`vendor` contents don't appear in results
