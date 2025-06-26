package pistis

import (
	"crypto/sha256"
	"errors"
	"fmt"
)

// Share represents a secret share owned by a party
type Share struct {
	Owner string
	Value []byte
}

// QueryResult holds the final query results, merkle proof, and secret shares
type QueryResult struct {
	Data         []tripleItem
	Proof        [][]byte
	SharedValues [][]Share
	QueryLog     []string
}

// simulateSPARQLQuery performs a two-step SPARQL-like query with join and verifiable proof generation.
func simulateSPARQLQuery(subjectTerm, predicate1, predicate2 string, t *trie, parties []string) (QueryResult, error) {
	queryLog := []string{fmt.Sprintf("Step 1: Match triples with subject = '%s'", subjectTerm)}

	// First triple pattern: ?s predicate1 ?o
	queryHash1 := hashStr1(subjectTerm)
	step1 := t.tripleQuery(queryHash1)
	if len(step1) == 0 {
		return QueryResult{}, errors.New("no result in step 1")
	}

	var joined []tripleItem
	var proofs [][]byte
	var shares [][]Share

	queryLog = append(queryLog, fmt.Sprintf("Found %d triples in step 1", len(step1)))

	// For each match in step 1, find join candidates in step 2
	for _, triple1 := range step1 {
		if !compareHash(triple1.Triple.predictHash, hashStr1(predicate1)) {
			continue
		}
		joinKey := triple1.Triple.objectHash
		queryLog = append(queryLog, fmt.Sprintf("Joining on object: %x", joinKey))

		step2 := t.tripleQuery(joinKey)
		for _, triple2 := range step2 {
			if compareHash(triple2.Triple.predictHash, hashStr1(predicate2)) {
				joined = append(joined, triple2)
				shares = append(shares, secretShare(triple2, parties))
				if proof := t.Root.generateProof(triple2.Triple.subjectHash); len(proof) > 0 {
					proofs = append(proofs, proof...)
				}
				queryLog = append(queryLog, fmt.Sprintf("Joined: %s", triple2.Data))
			}
		}
	}

	if len(joined) == 0 {
		return QueryResult{}, errors.New("join produced no result")
	}

	return QueryResult{
		Data:         joined,
		Proof:        proofs,
		SharedValues: shares,
		QueryLog:     queryLog,
	}, nil
}

// secretShare simulates secret sharing of RDF data
func secretShare(item tripleItem, parties []string) []Share {
	hash := sha256.Sum256([]byte(item.Data))
	shares := make([]Share, len(parties))
	for i, party := range parties {
		shares[i] = Share{
			Owner: party,
			Value: xorWithSalt(hash[:], byte(i+1)),
		}
	}
	return shares
}

func xorWithSalt(input []byte, salt byte) []byte {
	output := make([]byte, len(input))
	for i := range input {
		output[i] = input[i] ^ salt
	}
	return output
}

func compareHash(h1, h2 []byte) bool {
	if len(h1) != len(h2) {
		return false
	}
	for i := range h1 {
		if h1[i] != h2[i] {
			return false
		}
	}
	return true
}

func hashStr1(s string) []byte {
	h := sha256.Sum256([]byte(s))
	return h[:]
}

func debugQueryResult(res QueryResult) {
	fmt.Println("=== Query Execution Trace ===")
	for _, log := range res.QueryLog {
		fmt.Println(log)
	}
	fmt.Println("\n=== Results ===")
	for i, item := range res.Data {
		fmt.Printf("Result %d: %s\n", i+1, item.Data)
		fmt.Printf("Proof Size: %d, Shares: %d\n", len(res.Proof[i]), len(res.SharedValues[i]))
	}
}

// filterTripleByObject filters the result for additional condition like FILTER (?o = "Art")
func filterTripleByObject(items []tripleItem, objectTerm string) []tripleItem {
	var result []tripleItem
	objectHash := hashStr1(objectTerm)
	for _, item := range items {
		if compareHash(item.Triple.objectHash, objectHash) {
			result = append(result, item)
		}
	}
	return result
}

// extendedSPARQL supports filtering and query log for transparency
func extendedSPARQL(subject, pred1, pred2, filterObj string, t *trie, parties []string) (QueryResult, error) {
	res, err := simulateSPARQLQuery(subject, pred1, pred2, t, parties)
	if err != nil {
		return res, err
	}
	filtered := filterTripleByObject(res.Data, filterObj)
	res.Data = filtered
	res.QueryLog = append(res.QueryLog, fmt.Sprintf("Applied FILTER ?o = '%s' -> %d matches", filterObj, len(filtered)))
	return res, nil
}