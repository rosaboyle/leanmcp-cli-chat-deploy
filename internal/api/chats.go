package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListChats gets all chats for the authenticated user
func (c *Client) ListChats() ([]Chat, error) {
	resp, err := c.makeRequest("GET", "/api/chats", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list chats failed (status %d): %s", resp.StatusCode, string(body))
	}

	var chats []Chat
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		return nil, err
	}

	return chats, nil
}

// GetChat gets a specific chat by ID
func (c *Client) GetChat(chatID string) (*Chat, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/chats/id/%s", chatID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("chat not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get chat failed (status %d): %s", resp.StatusCode, string(body))
	}

	var chat Chat
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// GetChatHistory gets the message history for a chat
func (c *Client) GetChatHistory(chatID string) ([]ChatMessage, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/chats/id/%s/history/raw", chatID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("chat not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get chat history failed (status %d): %s", resp.StatusCode, string(body))
	}

	var messages []ChatMessage
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// CreateChat creates a new chat
func (c *Client) CreateChat(req CreateChatRequest) (*Chat, error) {
	resp, err := c.makeRequest("POST", "/api/chats", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create chat failed (status %d): %s", resp.StatusCode, string(body))
	}

	var chat Chat
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, err
	}

	return &chat, nil
}

// DeleteChat deletes a chat
func (c *Client) DeleteChat(chatID string) error {
	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/chats/id/%s", chatID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("chat not found")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete chat failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
