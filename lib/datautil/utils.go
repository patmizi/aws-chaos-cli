package datautil

func ListDiff(list1 []string, list2 []string) []string {
  var diff []string

  // Loop over the two lists comparing values
  for i := 0; i < 2; i++ {
    for _, l1 := range  list1 {
      found := false
      for _, l2 := range list2 {
        if l1 == l2 {
          found = true
          break
        }
      }
      // If string not found then add it to the diff list
      if !found {
        diff = append(diff, l1)
      }
    }
    // Swap lists on second iteration
    if i == 0 {
      list1, list2 = list2, list1
    }
  }

  return diff
}