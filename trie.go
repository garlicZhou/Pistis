package pistis

import (
	"bytes"
	"crypto/sha256"
	"github.com/syndtr/goleveldb/leveldb"
)

type triple struct {
	subjectHash []byte
	predictHash []byte
	objectHash  []byte
}

type tripleItem struct {
	Triple triple
	Data   string // optional payload
}

type node struct {
	parent    *node
	child     []*node
	childHash [][32]byte
	key       []byte
	value     []tripleItem
	hash      [32]byte
	isLeaf    bool
	isExtend  bool
}

type trie struct {
	Root     *node
	RootHash [32]byte
	DB       *leveldb.DB
}

func newTrie(db *leveldb.DB) *trie {
	return &trie{
		Root: &node{isExtend: false, isLeaf: false},
		DB:   db,
	}
}

func (n *node) updateHash(db *leveldb.DB) {
	if n == nil {
		return
	}

	h := sha256.New()
	n.childHash = nil // Reset before appending
	for _, child := range n.child {
		child.updateHash(db)
		h.Write(child.hash[:])
		n.childHash = append(n.childHash, child.hash) // Ensure sync
	}

	// Hash the current node's key and value
	h.Write(n.key)
	for _, v := range n.value {
		h.Write(v.Triple.subjectHash)
		h.Write(v.Triple.predictHash)
		h.Write(v.Triple.objectHash)
	}

	copy(n.hash[:], h.Sum(nil))

	// Optionally store in DB
	if db != nil && n.key != nil {
		db.Put(n.key, n.hash[:], nil)
	}
}

func (n *node) generateProof(target []byte) [][]byte {
	var proof [][]byte
	if n == nil {
		return nil
	}
	var dfs func(curr *node, path [][]byte) (bool, [][]byte)
	dfs = func(curr *node, path [][]byte) (bool, [][]byte) {
		if curr == nil {
			return false, path
		}
		for i, child := range curr.child {
			if child == nil {
				continue
			}
			if string(child.key) == string(target) {
				path = append(path, child.hash[:])
				return true, path
			}
			found, subPath := dfs(child, append(path, curr.childHash[i][:]))
			if found {
				return true, subPath
			}
		}
		return false, path
	}
	_, proof = dfs(n, proof)
	return proof
}

func (t *trie) tripleInsert(item tripleItem) {
	t.insertHash(item.Triple.subjectHash, item)
	t.insertHash(item.Triple.predictHash, item)
	t.insertHash(item.Triple.objectHash, item)
	t.Root.updateHash(t.DB)
	t.RootHash = t.Root.hash
}

func (t *trie) insertHash(word []byte, item tripleItem) {
	for _, j := range t.Root.child {
		if bytes.Equal(j.key, word) {
			j.value = append(j.value, item)
			j.updateHash(t.DB)
			return
		}
	}
	// If not found, create a new child node
	newNode := &node{
		key:    word,
		parent: t.Root,
		value:  []tripleItem{item},
		isLeaf: true,
	}
	newNode.updateHash(t.DB)
	t.Root.child = append(t.Root.child, newNode)
	t.Root.updateHash(t.DB)
}

func (n *node) nodeInsert(word []byte, item tripleItem, db *leveldb.DB) {
	for _, ch := range n.child {
		if bytes.Equal(ch.key, word) {
			ch.value = append(ch.value, item)
			ch.updateHash(db)
			return
		}
	}
	newNode := &node{
		key:    word,
		parent: n,
		value:  []tripleItem{item},
		isLeaf: true,
	}
	newNode.updateHash(db)

	n.child = append(n.child, newNode)
	n.updateHash(db)
}

func (t *trie) tripleQuery(hash []byte) []tripleItem {
	for _, ch := range t.Root.child {
		if len(ch.key) > 0 && hash[0] == ch.key[0] {
			return ch.nodeQuery(hash)
		}
	}
	return nil
}

func (n *node) nodeQuery(hash []byte) []tripleItem {
	if bytes.Equal(n.key, hash) {
		return n.value
	}
	for _, ch := range n.child {
		if bytes.Equal(ch.key, hash) {
			return ch.value
		}
	}
	return nil
}
