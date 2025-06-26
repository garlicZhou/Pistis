package pistis

import (
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func TestExtendedSPARQLQuery(t *testing.T) {
	dbPath := "./trie_vo_test.db"
	_ = os.RemoveAll(dbPath)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		t.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	tree := newTrie(db)
	parties := []string{"Alice", "Bob", "Carol"}

	triples := []tripleItem{
		{
			Triple: triple{
				subjectHash: hashStr("Bob"),
				predictHash: hashStr("owns"),
				objectHash:  hashStr("NFT123"),
			},
			Data: "Bob owns NFT123",
		},
		{
			Triple: triple{
				subjectHash: hashStr("NFT123"),
				predictHash: hashStr("type"),
				objectHash:  hashStr("Art"),
			},
			Data: "NFT123 is Art",
		},
		{
			Triple: triple{
				subjectHash: hashStr("NFT123"),
				predictHash: hashStr("creator"),
				objectHash:  hashStr("Alice"),
			},
			Data: "NFT123 created by Alice",
		},
	}

	for _, triple := range triples {
		tree.tripleInsert(triple)
	}

	res, err := extendedSPARQL("Bob", "owns", "type", "Art", tree, parties)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(res.Data) != 1 {
		t.Errorf("Expected 1 result for extended query, got %d", len(res.Data))
	}

	if len(res.SharedValues) == 0 || len(res.SharedValues[0]) != len(parties) {
		t.Errorf("Expected secret shares for all parties, got %d", len(res.SharedValues[0]))
	}

	if len(res.Proof) == 0 {
		t.Error("Expected Merkle proof, got none")
	}

	debugQueryResult(res)
}

func TestSPARQLJoinFailure(t *testing.T) {
	dbPath := "./trie_vo_test_empty.db"
	_ = os.RemoveAll(dbPath)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		t.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer db.Close()

	tree := newTrie(db)
	parties := []string{"OrgX", "OrgY"}

	tree.tripleInsert(tripleItem{
		Triple: triple{
			subjectHash: hashStr("Charlie"),
			predictHash: hashStr("owns"),
			objectHash:  hashStr("NFT456"),
		},
		Data: "Charlie owns NFT456",
	})

	res, err := extendedSPARQL("Charlie", "owns", "type", "Art", tree, parties)
	if err == nil {
		t.Error("Expected join failure due to missing triple, got success")
	}

	if len(res.Data) != 0 {
		t.Errorf("Expected 0 result for failed join, got %d", len(res.Data))
	}
}