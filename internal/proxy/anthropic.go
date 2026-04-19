package proxy

import "encoding/json"

// MessagesRequest mirrors the Anthropic Messages API request body.
// Content is RawMessage to handle both string and content-block array forms.
type MessagesRequest struct {
	Model     string            `json:"model"`
	Messages  []Message         `json:"messages"`
	MaxTokens int               `json:"max_tokens"`
	Stream    bool              `json:"stream"`
	System    string            `json:"system,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Message represents a single turn in the conversation.
// Content can be a plain string or a JSON array of content blocks.
type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// TextContent is the text variant of an Anthropic content block.
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ContentAsString returns the text of the message when content is a plain string.
// Returns empty string and false if content is not a plain string.
func (m *Message) ContentAsString() (string, bool) {
	var s string
	if err := json.Unmarshal(m.Content, &s); err == nil {
		return s, true
	}
	return "", false
}

// ContentAsBlocks returns the message content when it is an array of content blocks.
func (m *Message) ContentAsBlocks() ([]TextContent, bool) {
	var blocks []TextContent
	if err := json.Unmarshal(m.Content, &blocks); err == nil {
		return blocks, true
	}
	return nil, false
}

// SetContentString replaces the message content with a plain string.
func (m *Message) SetContentString(s string) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	m.Content = b
	return nil
}
