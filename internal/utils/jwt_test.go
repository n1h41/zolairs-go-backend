package utils

import "testing"

func TestGetUserId(t *testing.T) {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZXhwIjoxNTE2MjM5MDIyfQ.signature"

  userId, err := GetUserId(tokenString)
  if err != nil {
    t.Errorf("Expected no error, got %v", err)
  }

  if userId != "1234567890" {
    t.Errorf("Expected userId '1234567890', got %s", userId)
  }

}
