package main

import (
	"testing"
)

func TestNames(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		n, err := names(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != nil {
			t.Errorf("should return nil")
		}
	})

	t.Run("with []string", func(t *testing.T) {
		n, err := names([]interface{}{"hans"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n == nil {
			t.Errorf("should not return nil")
		}
		if len(n) != 1 || n[0] != "hans" {
			t.Errorf("unexpected return value: %v", n)
		}
	})

	t.Run("with string", func(t *testing.T) {
		n, err := names("hans")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n == nil {
			t.Errorf("should not return nil")
		}
		if len(n) != 1 || n[0] != "hans" {
			t.Errorf("unexpected return value: %v", n)
		}
	})

	t.Run("with map[string]string", func(t *testing.T) {
		_, err := names(map[string]string{})
		if err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("with []int", func(t *testing.T) {
		_, err := names([]interface{}{1, 2})
		if err == nil {
			t.Fatalf("expected an error")
		}
	})
}

func TestIds(t *testing.T) {
	t.Run("with nil", func(t *testing.T) {
		n, err := ids(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != nil {
			t.Errorf("should return nil")
		}
	})

	t.Run("with []int", func(t *testing.T) {
		n, err := ids([]interface{}{1})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n == nil {
			t.Errorf("should not return nil")
		}
		if len(n) != 1 || n[0] != 1 {
			t.Errorf("unexpected return value: %v", n)
		}
	})

	t.Run("with int", func(t *testing.T) {
		n, err := ids(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n == nil {
			t.Errorf("should not return nil")
		}
		if len(n) != 1 || n[0] != 1 {
			t.Errorf("unexpected return value: %v", n)
		}
	})

	t.Run("with map[string]string", func(t *testing.T) {
		_, err := ids(map[string]string{})
		if err == nil {
			t.Fatalf("expected an error")
		}
	})

	t.Run("with []string", func(t *testing.T) {
		_, err := ids([]interface{}{"a", "b"})
		if err == nil {
			t.Fatalf("expected an error")
		}
	})
}
