package main

import (
	"log"
	"path"

	"github.com/cipherowl-ai/addressdb/address"
	"github.com/cipherowl-ai/addressdb/store"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/piplabs/story-guardian/utils"
)

func main() {
	// number of addresses to generate
	count := 100
	// custom addresses to add
	customAddresses := []string{
		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"0x97DCA899a2278d010d678d64fBC7C718eD5D4939",
	}

	bf, err := store.NewBloomFilterStore(&address.EVMAddressHandler{})
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < count; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			log.Fatalf("Failed to generate key: %v", err)
		}
		addr := crypto.PubkeyToAddress(key.PublicKey).Hex()
		if err := bf.AddAddress(addr); err != nil {
			log.Fatalf("Failed to add address %s: %v", addr, err)
		}
	}

	for _, addr := range customAddresses {
		if err := bf.AddAddress(addr); err != nil {
			log.Fatalf("Failed to add address %s: %v", addr, err)
		}
	}

	filePath := path.Join(utils.GetDefaultPath(), "bloom_filter.gob")
	if err := bf.SaveToFile(filePath); err != nil {
		log.Fatalf("Failed to save Bloom filter to file: %v", err)
	}
}
