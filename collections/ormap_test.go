package collections

import (
	"testing"
)

func TestOrderedMap_AddItem(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.AddItem(1, "one")
	m.AddItem(2, "two")
	m.AddItem(3, "three")

	// Test if the items are added correctly
	if val, ok := m.GetItem(1); !ok || val != "one" {
		t.Errorf("Expected value 'one' for key 1, got '%s'", val)
	}
	if val, ok := m.GetItem(2); !ok || val != "two" {
		t.Errorf("Expected value 'two' for key 2, got '%s'", val)
	}
	if val, ok := m.GetItem(3); !ok || val != "three" {
		t.Errorf("Expected value 'three' for key 3, got '%s'", val)
	}

	// Test if existing item is updated correctly
	m.AddItem(2, "new two")
	if val, ok := m.GetItem(2); !ok || val != "new two" {
		t.Errorf("Expected value 'new two' for key 2, got '%s'", val)
	}
}

func TestOrderedMap_RemoveItem(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.AddItem(1, "one")
	m.AddItem(2, "two")
	m.AddItem(3, "three")

	// Test if item is removed correctly
	m.RemoveItem(2)
	if _, ok := m.GetItem(2); ok {
		t.Errorf("Expected item with key 2 to be removed")
	}

	// Test if other items are still accessible
	if val, ok := m.GetItem(1); !ok || val != "one" {
		t.Errorf("Expected value 'one' for key 1, got '%s'", val)
	}
	if val, ok := m.GetItem(3); !ok || val != "three" {
		t.Errorf("Expected value 'three' for key 3, got '%s'", val)
	}
}

func TestOrderedMap_GetAllItems(t *testing.T) {
	m := NewOrderedMap[int, string]()
	m.AddItem(1, "one")
	m.AddItem(2, "two")
	m.AddItem(3, "three")

	// Test if all items are returned correctly
	items := m.GetAllItems()
	if len(items) != 3 {
		t.Errorf("Expected 3 items, got %d", len(items))
	}
	if items[0].Key != 1 || items[0].Value != "one" {
		t.Errorf("Expected item 1 to be {1, 'one'}, got %+v", items[0])
	}
	if items[1].Key != 2 || items[1].Value != "two" {
		t.Errorf("Expected item 2 to be {2, 'two'}, got %+v", items[1])
	}
	if items[2].Key != 3 || items[2].Value != "three" {
		t.Errorf("Expected item 3 to be {3, 'three'}, got %+v", items[2])
	}
}
