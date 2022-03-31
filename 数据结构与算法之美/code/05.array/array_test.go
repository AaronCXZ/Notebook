package array

import "testing"

func TestInsert(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity-2; i++ {
		err := arr.Insert(uint(i), i+1)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	if err := arr.Insert(uint(6), 999); err != nil {
		t.Fatal(err.Error())
	}
	arr.Print()

	if err := arr.InsertToTail(777); err != nil {
		t.Fatal(err.Error())
	}
	arr.Print()
}

func TestDelete(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity; i++ {
		if err := arr.Insert(uint(i), i+1); err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	for i := 9; i >= 0; i-- {
		if _, err := arr.Delete(uint(i)); err != nil {
			t.Fatal(err.Error())
		}
		arr.Print()
	}
}

func TestFind(t *testing.T) {
	capacity := 10
	arr := NewArray(uint(capacity))
	for i := 0; i < capacity; i++ {
		if err := arr.Insert(uint(i), i+1); err != nil {
			t.Fatal(err.Error())
		}
	}
	arr.Print()

	t.Log(arr.Find(0))
	t.Log(arr.Find(6))
	t.Log(arr.Find(11))
}
