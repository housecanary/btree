package btree

import (
	"encoding/binary"
	"fmt"
	"io"
)

func (t *BTree) Save(
	f io.Writer,
	saveItem func (w io.Writer, value Item) error,
) (err error) {
	if err = binary.Write(f, binary.BigEndian, uint64(t.degree)); err != nil {
		return
	}
	fmt.Printf("Wrote degree\n")

	if err = binary.Write(f, binary.BigEndian, uint64(t.length)); err != nil {
		return
	}
	fmt.Printf("Wrote length\n")

	gotTree := t.root != nil
	if err = binary.Write(f, binary.BigEndian, gotTree); err != nil {
		return
	}
	fmt.Printf("Wrote gotTree\n")

	if t.root != nil {
		if err = t.root.save(f, saveItem); err != nil {
			return
		}
	}

	return
}

func (n *node) save(
	f io.Writer,
	saveItem func (w io.Writer, item Item) error,
) (err error) {
	nItems := len(n.items)
	if err = binary.Write(f, binary.BigEndian, uint8(nItems)); err != nil {
		return
	}
	fmt.Printf("Wrote nItems: %v\n", nItems)

	gotChildren := len(n.children) > 0
	if err = binary.Write(f, binary.BigEndian, gotChildren); err != nil {
		return
	}
	fmt.Printf("Wrote gotChildren: %v\n", gotChildren)
	// values on this node
	for i := 0; i < nItems; i++ {
		item := n.items[i]
		if err = saveItem(f, item); err != nil {
			return
		}
	}
	// children
	if gotChildren {
		for i := 0; i <= nItems; i++ {
			if err = n.children[i].save(f, saveItem); err != nil {
				return
			}
		}
	}

	return
}

func Load(
	f io.Reader,
	loadItem func (r io.Reader, obuf []byte) (Item, []byte, error),
) (t *BTree, err error) {
	t = &BTree{}
	var word uint64

	if err = binary.Read(f, binary.BigEndian, &word); err != nil {
		return
	}
	fmt.Printf("Read degree\n")
	t.degree = int(word)

	if err = binary.Read(f, binary.BigEndian, &word); err != nil {
		return
	}
	fmt.Printf("Read length\n")
	t.length = int(word)

	var gotTree bool
	if err = binary.Read(f, binary.BigEndian, &gotTree); err != nil {
		return
	}
	fmt.Printf("Read gotTree: %v\n", gotTree)

	if gotTree {
		// this buffer will be re-used or replaced for a larger one, as needed
		buf := make([]byte, 0)
		if t.root, buf, err = load(f, buf, loadItem); err != nil {
			return
		}
	}
	return
}

func load(
	f io.Reader,
	oldBuf []byte,
	loadItem func (r io.Reader, obuf []byte) (Item, []byte, error),
) (n *node, buf []byte, err error) {
	buf = oldBuf[:]
	n = &node{}

	var short uint8
	if err = binary.Read(f, binary.BigEndian, &short); err != nil {
		return
	}
	fmt.Printf("Read numItems: %d\n", short)
	nItems := int(short)

	var gotChildren bool
	if err = binary.Read(f, binary.BigEndian, &gotChildren); err != nil {
		return
	}
	fmt.Printf("Read gotChildren: %v\n", gotChildren)

	// values on this node
	var item Item
	n.items = make([]Item, nItems)
	for i := 0; i < nItems; i++ {
		if item, buf, err = loadItem(f, buf); err != nil {
			return
		}
		n.items[i] = item
	}
	// children
	if gotChildren {
		n.children = make([]*node, nItems+1)
		for i := 0; i <= nItems; i++ {
			if n.children[i], buf, err = load(f, buf, loadItem); err != nil {
				return
			}
		}
	}

	return
}
