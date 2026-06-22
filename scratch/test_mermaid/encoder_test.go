package mermaid

import (
	"strings"
	"testing"
)

func TestGenerateLiveURL(t *testing.T) {
	code := `graph TD
    A[Start] --> B[Process]
    B --> C[End]`

	url, err := GenerateLiveURL(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prefix := "https://mermaid.live/edit#pako:"
	if !strings.HasPrefix(url, prefix) {
		t.Fatalf("expected URL to start with %q, got %q", prefix, url)
	}
}

func TestInsertLiveLinks(t *testing.T) {
	input := `Some introductory text.

```mermaid
graph TD
    A --> B
```

Some concluding text.`

	output := InsertLiveLinks(input)
	if !strings.Contains(output, "[Open in Mermaid Live Editor ↗]") {
		t.Errorf("expected output to contain live editor link, got:\n%s", output)
	}

	// Verify it does not affect non-mermaid blocks.
	nonMermaidInput := ````go
package main
import "fmt"
func main() {
    fmt.Println("hello")
}
````
`
	nonMermaidOutput := InsertLiveLinks(nonMermaidInput)
	if strings.Contains(nonMermaidOutput, "Mermaid Live Editor") {
		t.Errorf("did not expect live editor link for Go code block, got:\n%s", nonMermaidOutput)
	}
}
