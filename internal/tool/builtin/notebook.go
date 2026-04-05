package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/strings77wzq/claude-code-Go/internal/permission"
	"github.com/strings77wzq/claude-code-Go/internal/tool"
)

// Notebook represents a Jupyter notebook (nbformat v4)
type Notebook struct {
	NBFormat   int             `json:"nbformat"`
	NBFormatMA int             `json:"nbformat_minor"`
	Metadata   json.RawMessage `json:"metadata"`
	Cells      []Cell          `json:"cells"`
}

// Cell represents a notebook cell
type Cell struct {
	CellType       string          `json:"cell_type"`
	Source         interface{}     `json:"source"`
	Metadata       json.RawMessage `json:"metadata"`
	Outputs        json.RawMessage `json:"outputs,omitempty"`
	ExecutionCount interface{}     `json:"execution_count,omitempty"`
}

// NotebookTool provides operations for editing Jupyter notebook files
type NotebookTool struct {
	workingDir string
}

// NewNotebookTool creates a new NotebookTool instance
func NewNotebookTool(workingDir string) tool.Tool {
	return &NotebookTool{
		workingDir: workingDir,
	}
}

func (n *NotebookTool) Name() string {
	return "NotebookEdit"
}

func (n *NotebookTool) Description() string {
	return "Edit Jupyter notebook (.ipynb) files"
}

func (n *NotebookTool) InputSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "Path to the Jupyter notebook file (.ipynb)",
			},
			"operation": map[string]any{
				"type":        "string",
				"description": "Operation to perform: read, edit, add_cell, delete_cell, execute",
				"enum":        []string{"read", "edit", "add_cell", "delete_cell", "execute"},
			},
			"cell_index": map[string]any{
				"type":        "number",
				"description": "Index of the cell to edit, delete, or execute (0-based)",
			},
			"source": map[string]any{
				"type":        "string",
				"description": "Source code or markdown content for edit/add_cell operations",
			},
			"language": map[string]any{
				"type":        "string",
				"description": "Language for new cells: python (default) or javascript",
				"enum":        []string{"python", "javascript"},
			},
			"cell_type": map[string]any{
				"type":        "string",
				"description": "Cell type for add_cell: code (default) or markdown",
				"enum":        []string{"code", "markdown"},
			},
		},
		"required": []string{"file_path", "operation"},
	}
}

func (n *NotebookTool) RequiresPermission() bool {
	return true
}

func (n *NotebookTool) RequiredPermissionLevel() permission.PermissionLevel {
	return permission.LevelWorkspaceWrite
}

// getStringSlice converts source interface{} to []string
func getStringSlice(source interface{}) ([]string, bool) {
	switch s := source.(type) {
	case string:
		return []string{s}, true
	case []interface{}:
		result := make([]string, len(s))
		for i, item := range s {
			if str, ok := item.(string); ok {
				result[i] = str
			} else {
				return nil, false
			}
		}
		return result, true
	default:
		return nil, false
	}
}

// setStringSlice converts []string to the appropriate JSON source format
func setStringSlice(lines []string) interface{} {
	if len(lines) == 1 {
		return lines[0]
	}
	var result []string
	for _, line := range lines {
		result = append(result, line)
	}
	return result
}

func (n *NotebookTool) Execute(ctx context.Context, input map[string]any) tool.Result {
	filePath, ok := input["file_path"].(string)
	if !ok || filePath == "" {
		return tool.Error("file_path is required")
	}

	if err := ValidatePath(filePath, n.workingDir); err != nil {
		return tool.Error(err.Error())
	}

	operation, ok := input["operation"].(string)
	if !ok || operation == "" {
		return tool.Error("operation is required")
	}

	switch operation {
	case "read":
		return n.readNotebook(filePath)
	case "edit":
		return n.editCell(filePath, input)
	case "add_cell":
		return n.addCell(filePath, input)
	case "delete_cell":
		return n.deleteCell(filePath, input)
	case "execute":
		return n.executeCell(filePath, input)
	default:
		return tool.Error(fmt.Sprintf("unknown operation: %s", operation))
	}
}

// readNotebook reads and returns the notebook content
func (n *NotebookTool) readNotebook(filePath string) tool.Result {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tool.Error(fmt.Sprintf("file not found: %s", filePath))
		}
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	var nb Notebook
	if err := json.Unmarshal(content, &nb); err != nil {
		return tool.Error(fmt.Sprintf("failed to parse notebook: %v", err))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Notebook: %s\n", filePath))
	sb.WriteString(fmt.Sprintf("Format: %d.%d\n", nb.NBFormat, nb.NBFormatMA))
	sb.WriteString(fmt.Sprintf("Cells: %d\n\n", len(nb.Cells)))

	for i, cell := range nb.Cells {
		sb.WriteString(fmt.Sprintf("--- Cell %d (%s) ---\n", i, cell.CellType))

		lines, ok := getStringSlice(cell.Source)
		if !ok {
			sb.WriteString("(unable to parse source)\n")
		} else {
			for _, line := range lines {
				sb.WriteString(line + "\n")
			}
		}

		if cell.CellType == "code" && len(cell.Outputs) > 0 {
			sb.WriteString("\n[Outputs present]\n")
		}
		sb.WriteString("\n")
	}

	return tool.Success(sb.String())
}

// editCell modifies a cell's source at the given index
func (n *NotebookTool) editCell(filePath string, input map[string]any) tool.Result {
	cellIndex, ok := input["cell_index"]
	if !ok {
		return tool.Error("cell_index is required for edit operation")
	}

	idx, err := toInt(cellIndex)
	if err != nil {
		return tool.Error("cell_index must be a number")
	}

	source, ok := input["source"].(string)
	if !ok || source == "" {
		return tool.Error("source is required for edit operation")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	var nb Notebook
	if err := json.Unmarshal(content, &nb); err != nil {
		return tool.Error(fmt.Sprintf("failed to parse notebook: %v", err))
	}

	if idx < 0 || idx >= len(nb.Cells) {
		return tool.Error(fmt.Sprintf("cell_index %d out of range (0-%d)", idx, len(nb.Cells)-1))
	}

	nb.Cells[idx].Source = setStringSlice([]string{source})

	newContent, err := json.MarshalIndent(nb, "", "  ")
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to serialize notebook: %v", err))
	}

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return tool.Error(fmt.Sprintf("failed to write file: %v", err))
	}

	return tool.Success(fmt.Sprintf("Successfully edited cell %d in %s", idx, filePath))
}

// addCell inserts a new cell at the given index
func (n *NotebookTool) addCell(filePath string, input map[string]any) tool.Result {
	cellIndex, ok := input["cell_index"]

	var idx int
	var err error
	if ok {
		idx, err = toInt(cellIndex)
		if err != nil {
			return tool.Error("cell_index must be a number")
		}
	}

	source := ""
	if s, ok := input["source"].(string); ok {
		source = s
	}

	cellType := "code"
	if ct, ok := input["cell_type"].(string); ok {
		cellType = ct
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	var nb Notebook
	if err := json.Unmarshal(content, &nb); err != nil {
		return tool.Error(fmt.Sprintf("failed to parse notebook: %v", err))
	}

	newCell := Cell{
		CellType: cellType,
		Source:   setStringSlice([]string{source}),
		Metadata: json.RawMessage(`{}`),
	}

	if cellType == "code" {
		newCell.Outputs = json.RawMessage(`[]`)
		newCell.ExecutionCount = json.RawMessage("null")
	}

	// Insert at the specified index or append
	if ok && idx >= 0 && idx <= len(nb.Cells) {
		nb.Cells = append(nb.Cells[:idx], append([]Cell{newCell}, nb.Cells[idx:]...)...)
	} else {
		nb.Cells = append(nb.Cells, newCell)
		idx = len(nb.Cells) - 1
	}

	newContent, err := json.MarshalIndent(nb, "", "  ")
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to serialize notebook: %v", err))
	}

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return tool.Error(fmt.Sprintf("failed to write file: %v", err))
	}

	return tool.Success(fmt.Sprintf("Successfully added %s cell at index %d in %s", cellType, idx, filePath))
}

// deleteCell removes a cell at the given index
func (n *NotebookTool) deleteCell(filePath string, input map[string]any) tool.Result {
	cellIndex, ok := input["cell_index"]
	if !ok {
		return tool.Error("cell_index is required for delete_cell operation")
	}

	idx, err := toInt(cellIndex)
	if err != nil {
		return tool.Error("cell_index must be a number")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	var nb Notebook
	if err := json.Unmarshal(content, &nb); err != nil {
		return tool.Error(fmt.Sprintf("failed to parse notebook: %v", err))
	}

	if idx < 0 || idx >= len(nb.Cells) {
		return tool.Error(fmt.Sprintf("cell_index %d out of range (0-%d)", idx, len(nb.Cells)-1))
	}

	nb.Cells = append(nb.Cells[:idx], nb.Cells[idx+1:]...)

	newContent, err := json.MarshalIndent(nb, "", "  ")
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to serialize notebook: %v", err))
	}

	if err := os.WriteFile(filePath, newContent, 0644); err != nil {
		return tool.Error(fmt.Sprintf("failed to write file: %v", err))
	}

	return tool.Success(fmt.Sprintf("Successfully deleted cell %d from %s", idx, filePath))
}

// executeCell runs the code in a cell and returns the output
func (n *NotebookTool) executeCell(filePath string, input map[string]any) tool.Result {
	cellIndex, ok := input["cell_index"]
	if !ok {
		return tool.Error("cell_index is required for execute operation")
	}

	idx, err := toInt(cellIndex)
	if err != nil {
		return tool.Error("cell_index must be a number")
	}

	language := "python"
	if lang, ok := input["language"].(string); ok {
		language = lang
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return tool.Error(fmt.Sprintf("failed to read file: %v", err))
	}

	var nb Notebook
	if err := json.Unmarshal(content, &nb); err != nil {
		return tool.Error(fmt.Sprintf("failed to parse notebook: %v", err))
	}

	if idx < 0 || idx >= len(nb.Cells) {
		return tool.Error(fmt.Sprintf("cell_index %d out of range (0-%d)", idx, len(nb.Cells)-1))
	}

	cell := nb.Cells[idx]
	if cell.CellType != "code" {
		return tool.Error(fmt.Sprintf("cannot execute markdown cell"))
	}

	lines, ok := getStringSlice(cell.Source)
	if !ok || len(lines) == 0 {
		return tool.Error("cell has no source code to execute")
	}

	code := ""
	for _, line := range lines {
		code += line + "\n"
	}

	var output []byte
	switch language {
	case "python":
		cmd := exec.Command("python", "-c", code)
		output, err = cmd.CombinedOutput()
	case "javascript":
		cmd := exec.Command("node", "-e", code)
		output, err = cmd.CombinedOutput()
	default:
		return tool.Error(fmt.Sprintf("unsupported language: %s", language))
	}

	if err != nil {
		return tool.Error(fmt.Sprintf("execution error: %v\n%s", err, string(output)))
	}

	return tool.Success(string(output))
}

// toInt converts various numeric types to int
func toInt(v interface{}) (int, error) {
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int64:
		return int(n), nil
	case int:
		return n, nil
	case string:
		return strconv.Atoi(n)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v)
	}
}

// toInt converts various numeric types to int
