package pistis

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"math/rand"
	"time"
)

type encryptedNode struct {
	EncryptedData []byte
	Hash          [32]byte
}

func pad(data []byte, blockSize int) []byte {
	padLen := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func aesEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data = pad(data, block.BlockSize())
	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	mode.CryptBlocks(ciphertext, data)
	return ciphertext, nil
}

func prfKey(seed []byte) []byte {
	h := sha256.Sum256(seed)
	return h[:16] // AES-128 key
}

func serializeNode(n *node) []byte {
	var buf bytes.Buffer
	buf.Write(n.key)
	for _, v := range n.value {
		buf.Write(v.Triple.subjectHash)
		buf.Write(v.Triple.predictHash)
		buf.Write(v.Triple.objectHash)
		buf.WriteString(v.Data)
	}
	return buf.Bytes()
}

func encryptTrieNodes(n *node, seed []byte) ([]encryptedNode, error) {
	rand.Seed(time.Now().UnixNano())
	var encrypted []encryptedNode
	queue := []*node{n}
	permute := rand.Perm(len(n.child) + 1) // +1 for root
	key := prfKey(seed)

	for _, idx := range permute {
		if idx == 0 {
			enc, err := aesEncrypt(serializeNode(n), key)
			if err != nil {
				return nil, err
			}
			encrypted = append(encrypted, encryptedNode{EncryptedData: enc, Hash: n.hash})
			queue = append(queue, n.child...)
			continue
		}
		if idx-1 < len(n.child) {
			c := n.child[idx-1]
			enc, err := aesEncrypt(serializeNode(c), key)
			if err != nil {
				return nil, err
			}
			encrypted = append(encrypted, encryptedNode{EncryptedData: enc, Hash: c.hash})
			queue = append(queue, c.child...)
		}
	}

	return encrypted, nil
}

func decryptNode(enc []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(enc)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	plaintext := make([]byte, len(enc))
	mode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	mode.CryptBlocks(plaintext, enc)
	return plaintext, nil
}
