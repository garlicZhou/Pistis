
# Pistis: Encrypted Merkle Trie and Secure SPARQL Verification

This project is an implementation of the system described in the paper:

**"Pistis: A Decentralized Knowledge Graph Platform Enabling Ownership-Preserving SPARQL Querying" (VLDB 2025)**

It consists of three main components:

1. **Merkle Trie Construction (`trie.go`)**: A prefix tree for RDF triples, similar to a Merkle Patricia Trie.
2. **Encrypted Merkle Trie (`emst.go`)**: Implements the STRP algorithm, including padding, encryption, and permutation.
3. **VO-SPARQL Secure Query Verification (`vo_sparql.go`)**: Performs secure SPARQL query verification using secret sharing and MPC principles.

---

## 🛠 Project Structure

```bash
.
├── trie.go         # Implements the Merkle Trie for RDF data
├── emst.go         # Encrypted Merkle Trie with STRP: Padding, Encrypting, Permuting
├── vo_sparql.go    # Secure SPARQL verification using MPC (ABY simulated)
├── go.mod          # Go module definition
├── README.md       # This file
```

---

## 🧱 Requirements

- Go 1.17+
- LevelDB (for trie persistence)
- [ABY Framework](https://encrypto.de/code/ABY) (if using real MPC, optional)
- Unix-like OS (recommended for subprocess execution)

Install Go and dependencies:

```bash
# Ubuntu
sudo apt install golang libleveldb-dev

# MacOS (with Homebrew)
brew install go leveldb
```

---

## ▶️ How to Run

### 1. Clone the repository

```bash
git clone https://github.com/your-org/pistis
cd pistis
```

### 2. Initialize Go module

```bash
go mod init pistis
go mod tidy
```

### 3. Run the project

```bash
go run trie_test.go emst_test.go vo_sparql_test.go
```

---

## 🔐 Merkle Trie Logic (trie.go)

The trie supports:
- Insertion of RDF triples with S/P/O hashing
- RLP encoding and SHA-256 for Merkle hash generation
- On-disk persistence using LevelDB

---

## 🔒 EMST: Encrypted Merkle Semantic Trie (emst.go)

Implements the **STRP** Algorithm:

1. **Padding** – Ensures data aligns with AES block size
2. **Encrypting** – Uses AES-256 CBC encryption with a random IV
3. **Permuting** – Randomly shuffles node orders to obfuscate access pattern

---

## 🔎 Secure VO-SPARQL (vo_sparql.go)

- Simulates **secret sharing** of RDF query answers
- Demonstrates verifiable computation with reconstruction

---

---

## 🔬 Integration with Real Blockchain Systems

To test the performance of **Pistis** in real blockchain scenarios, please replace the `trie.go` file in the **ETHMST** project with the `trie.go` file from this project.

Also, add the other Go files (`emst.go`, `vo_sparql.go`, etc.) to the appropriate directory in ETHMST for compatibility.

The ETHMST repository address is:  
👉 [ETHMST GitHub Repository](https://github.com/garlicZhou/ETHMST)

## 📂 License

This is a research prototype inspired by the VLDB 2025 paper.

For academic use only.
