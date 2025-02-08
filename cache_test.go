package cache

import "testing"

func TestGetNotFound(t *testing.T) {
	c := New(1)
	if _, ok := c.Get("a"); ok != false {
		t.Errorf("c.Get('a') = %t, want false", ok)
	}
}

func TestGetFound(t *testing.T) {
	c := New(1)
	c.Set("a", 1)
	if val, ok := c.Get("a"); val != 1 || ok != true {
		t.Errorf("c.Get('a') = %d, %t, want 1, true", val, ok)
	}
}

func TestSetExistingKey(t *testing.T) {
	c := New(1)
	c.Set("a", 1)
	c.Set("a", 2)
	if val, ok := c.Get("a"); val != 2 || ok != true {
		t.Errorf("c.Get('a') = %d, %t, want 2, true", val, ok)
	}
}

func TestEvictLeastFrequentlyUsedKey(t *testing.T) {
	c := New(2)
	c.Set("a", 1)
	// set b's freq to 3
	c.Set("b", 2)
	c.Get("b")
	c.Get("b")

	// set a to be more recently accessed, but has freq of 2
	c.Get("a")

	// insert new key, expect a to be evicted, and b is kept intact, c is
	// inserted
	c.Set("c", 1)

	if _, ok := c.Get("a"); ok != false {
		t.Errorf("c.Get('a') = _, %t, want _, false", ok)
	}
	if val, ok := c.Get("b"); val != 2 || ok != true {
		t.Errorf("c.Get('b') = %d, %t, want 2, true", val, ok)
	}
	if val, ok := c.Get("c"); val != 1 || ok != true {
		t.Errorf("c.Get('b') = %d, %t, want 1, true", val, ok)
	}
}

func TestEvictLeastRecentlyUsedKeyIfSameFrequency(t *testing.T) {
	c := New(2)
	c.Set("a", 2)
	// set b's freq to 2
	c.Set("b", 2)
	c.Get("b")

	// set a to be more recently accessed has freq of 2 (same as b)
	c.Get("a")

	// insert new key, expect b to be evicted, a is kept intact, c is inserted
	c.Set("c", 1)

	if _, ok := c.Get("b"); ok != false {
		t.Errorf("c.Get('b') = _, %t, want _, false", ok)
	}
	if val, ok := c.Get("a"); val != 2 || ok != true {
		t.Errorf("c.Get('c') = %d, %t, want 2, true", val, ok)
	}
	if val, ok := c.Get("c"); val != 1 || ok != true {
		t.Errorf("c.Get('b') = %d, %t, want 1, true", val, ok)
	}
}
