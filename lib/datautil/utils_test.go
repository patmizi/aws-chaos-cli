package datautil

import (
  "reflect"
  "testing"
)

func TestListDiff(t *testing.T) {
  l1 := []string{"a", "b", "c"}
  l2 := []string{"b", "c", "d"}
  diff := ListDiff(l1, l2)
  actual := []string{"a", "d"}

  if !reflect.DeepEqual(diff, actual) {
    t.Errorf("got %v want %v", diff, actual)
  }
}