package btree

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

type itemT uint32

func (item *itemT) Less(other Item, ctx interface{}) bool {
	return *item < *other.(*itemT)
}

func itemSaver(w io.Writer, itm Item) (err error) {
	item := itm.(*itemT)
	if err = binary.Write(w, binary.BigEndian, item); err != nil {
		return
	}
	return
}

func itemLoader(r io.Reader, obuf []byte) (item Item, buf []byte, err error) {
	var itm itemT
	if err = binary.Read(r, binary.BigEndian, &itm); err != nil {
		return
	}
	return &itm, buf,nil
}

func TestSaveLoadBTree(t *testing.T) {
	tr := New(*btreeDegree, nil)

	for _, i := range rand.Perm(256) {
		item := itemT(i)
		if x := tr.ReplaceOrInsert(&item); x != nil {
			t.Fatal("insert found item", item)
		}
	}

	var f *os.File
	var err error
	fileName := "/tmp/btree_save"
	f, err = os.Create(fileName)
	if err != nil {
		t.Fatal("creating failed")
	}

	if err = tr.Save(f, itemSaver); err != nil {
		t.Fatal("saving failed")
	}
	if f.Close() != nil {
		t.Fatal("closing failed")
	}

	f, err = os.Open(fileName)
	if err != nil {
		t.Fatal("opening failed")
	}

	var newTr *BTree
	if newTr, err = Load(f, itemLoader); err != nil {
		t.Fatal("loading failed")
	}

	if f.Close() != nil {
		t.Fatal("closing failed")
	}

	fmt.Printf("Orig tree: degree %d length %d\n", tr.degree, tr.length)
	fmt.Printf("New tree: degree %d length %d\n", newTr.degree, newTr.length)

	fmt.Printf("Old tree: %v\n", tr.root)
	fmt.Printf("New tree: %v\n", newTr.root)

	oldTree := all(tr)
	newTree := all(newTr)
	if !reflect.DeepEqual(oldTree, newTree) {
		t.Fatalf("mismatch:\n old: %v\nnew: %v", oldTree, newTree)
	}
}
