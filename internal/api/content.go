package api

// ExtractTextContent extracts all text from content blocks.
func ExtractTextContent(blocks []ContentBlock) string {
	var text string
	for _, block := range blocks {
		if block.Type == "text" {
			text += block.Text
		}
	}
	return text
}
