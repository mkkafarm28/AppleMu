package main

import (
	"fmt"
	"os"

	"amdecrypt/pkg/decrypt"
)

func main() {
	if len(os.Args) != 7 {
		fmt.Fprintf(os.Stderr, "usage: %s <agentIp> <mp4decryptPath> <id> <key> <inputPath> <outputPath>\n", os.Args[0])
		os.Exit(1)
	}

	agentIp := os.Args[1]
	mp4decryptPath := os.Args[2]
	id := os.Args[3]
	key := os.Args[4]
	inputPath := os.Args[5]
	outputPath := os.Args[6]

	info, err := decrypt.ExtractSong(inputPath)
	if err != nil {
		panic(err)
	}

	keys := []string{decrypt.PrefetchKey, key}
	err = decrypt.DecryptSong(agentIp, mp4decryptPath, outputPath, id, info, keys)
	if err != nil {
		panic(err)
	}
}
