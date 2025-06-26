package pistis

import (
	"crypto/sha256"
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func hashStr(s string) []byte {
	h := sha256.Sum256([]byte(s))
	return h[:]
}

func TestTrieInsertAndQuery(t *testing.T) {
	dbPath := "./trie_test_db"
	_ = os.RemoveAll(dbPath)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		t.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	tr := newTrie(db)

	t1 := tripleItem{
		Triple: triple{
			subjectHash: hashStr("Alice"),
			predictHash: hashStr("knows"),
			objectHash:  hashStr("Bob"),
		},
		Data: "Alice knows Bob",
	}

	t2 := tripleItem{
		Triple: triple{
			subjectHash: hashStr("Bob"),
			predictHash: hashStr("likes"),
			objectHash:  hashStr("Pizza"),
		},
		Data: "Bob likes Pizza",
	}

	t3 := tripleItem{
		Triple: triple{
			subjectHash: hashStr("Alice"),
			predictHash: hashStr("likes"),
			objectHash:  hashStr("Music"),
		},
		Data: "Alice likes Music",
	}

	tr.tripleInsert(t1)
	tr.tripleInsert(t2)
	tr.tripleInsert(t3)

	if got := tr.tripleQuery(hashStr("Alice")); len(got) != 2 {
		t.Errorf("Expected 2 results for Alice, got %d", len(got))
	}

	if got := tr.tripleQuery(hashStr("likes")); len(got) != 2 {
		t.Errorf("Expected 2 results for likes, got %d", len(got))
	}

	if got := tr.tripleQuery(hashStr("Pizza")); len(got) != 1 || got[0].Data != "Bob likes Pizza" {
		t.Errorf("Expected Bob likes Pizza, got %+v", got)
	}
}
