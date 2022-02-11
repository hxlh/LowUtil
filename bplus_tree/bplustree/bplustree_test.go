package bplustree

import (
	"testing"
)

func TestDelete1(t *testing.T) {
	b := &BPlusTree{}
	b.Insert(100)
	b.Insert(20)
	b.Insert(110)

	b.Insert(1)
	b.Display("i 1")
	b.Insert(200)
	b.Display("i 200")
	b.Insert(300)
	b.Display("i 300")
	b.Insert(400)
	b.Display("i 400")
	b.Insert(500)
	b.Display("i 500")
	b.Insert(600)
	b.Display("i 600")
	b.Insert(700)
	b.Display("i 700")
	b.Insert(220)
	b.Display("insert")

	b.Delete(110)
	b.Display("d 110")
	b.Delete(100)
	b.Display("d 100")

	b.Delete(300)
	b.Display("d 300")

	b.Delete(500)
	b.Display("d 500")

	b.Insert(800)
	b.Display("i 800")
	b.Insert(900)
	b.Display("i 900")

	b.Delete(700)
	b.Display("d 900")
	b.Delete(200)
	b.Display("d 200")
	b.Delete(220)
	b.Display("d 220")

	b.Delete(400)
	b.Display("d 400")
	b.Delete(600)
	b.Display("d 600")

	b.Delete(900)
	b.Display("d 900")

	b.Delete(800)
	b.Display("d 800")

	b.Delete(1)
	b.Display("d 1")

	

	b.Delete(20)
	b.Display("d 20")

	b.Insert(1)
	b.Display("d 1")

}
