package pistis

import (
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestEncryptedTrie(t *testing.T) {
	dbPath := "./trie_enc_test_db"
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
			subjectHash: hashStr("Alice"),
			predictHash: hashStr("likes"),
			objectHash:  hashStr("Music"),
		},
		Data: "Alice likes Music",
	}

	tr.tripleInsert(t1)
	tr.tripleInsert(t2)

	seed := hashStr("encryption-seed")
	encryptedNodes, err := encryptTrieNodes(tr.Root, seed)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	if len(encryptedNodes) == 0 {
		t.Fatal("Expected encrypted nodes, got none")
	}

	decrypted, err := decryptNode(encryptedNodes[0].EncryptedData, prfKey(seed))
	if err != nil {
		t.Errorf("Decryption failed: %v", err)
	} else if len(decrypted) == 0 {
		t.Error("Expected non-empty decrypted data")
	}
}


