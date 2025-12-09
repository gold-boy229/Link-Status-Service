package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"Link-Status-Service/internal/utils"
)

type linkRepository struct {
	counter           *atomic.Uint64
	mapHashToLinkNum  sync.Map
	mapLinkNumToLinks sync.Map
}

func NewLinkRepository() *linkRepository {
	repo := linkRepository{counter: &atomic.Uint64{}}
	repo.counter.Add(1)
	return &repo
}

func (repo *linkRepository) GetLinkNum(ctx context.Context, links []string) (linkNum int, isNew bool, err error) {
	sortedLinks := utils.SortStrings(links)
	hash64 := getHashOfLinks(sortedLinks)
	linkNum, isNew = repo.hashToLinkNumGetOrCreate(hash64)
	return linkNum, isNew, nil
}

func getHashOfLinks(links []string) uint64 {
	sb := bytes.NewBufferString("")
	for _, link := range links {
		sb.WriteString(link)
	}

	hash := fnv.New64a()
	hash.Write(sb.Bytes())
	return hash.Sum64()
}

func (repo *linkRepository) hashToLinkNumGetOrCreate(hash uint64) (linkNum int, isNew bool) {
	newLinkNum := repo.counter.Load()
	actualLinkNum, alreadyExisted := repo.mapHashToLinkNum.LoadOrStore(hash, int(newLinkNum))

	isNew = !alreadyExisted
	if isNew {
		repo.counter.Add(1)
	}

	linkNum, _ = actualLinkNum.(int)
	return linkNum, isNew
}

func (repo *linkRepository) StoreLinks(ctx context.Context, links []string, linkNum int) error {
	sortedLinks := utils.SortStrings(links)
	repo.mapLinkNumToLinks.Store(linkNum, sortedLinks)

	return nil
}

func (repo *linkRepository) GetLinksByLinkNum(ctx context.Context, linkNum int) ([]string, error) {
	linksAny, exists := repo.mapLinkNumToLinks.Load(linkNum)
	if !exists {
		return []string{}, nil
	}

	links, ok := linksAny.([]string)
	if !ok {
		return []string{}, errors.New("cannot convert links to []string")
	}

	return links, nil
}

func (repo *linkRepository) StoreDataToJSON() error {
	// Prepare and Save mapHashToLinkNum
	hashToLinkNumMap := make(map[uint64]int)
	repo.mapHashToLinkNum.Range(func(key, value any) bool {
		hash, ok := key.(uint64)
		if !ok {
			return false
		}
		linkNum, ok := value.(int)
		if !ok {
			return false
		}
		hashToLinkNumMap[hash] = linkNum
		return true
	})
	if err := saveMapToFile("./data/hash_to_link_num.json", hashToLinkNumMap); err != nil {
		return err
	}

	// Prepare and Save mapLinkNumToLinks
	linkNumToLinksMap := make(map[int][]string)
	repo.mapLinkNumToLinks.Range(func(key, value any) bool {
		linkNum, ok := key.(int)
		if !ok {
			return false
		}
		links, ok := value.([]string)
		if !ok {
			return false
		}
		linkNumToLinksMap[linkNum] = links
		return true
	})
	if err := saveMapToFile("./data/link_num_to_links.json", linkNumToLinksMap); err != nil {
		return err
	}

	return nil
}

func saveMapToFile(filename string, data interface{}) error {
	// Ensure the directory exists
	dirName := filepath.Dir(filename)
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirName, err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data for %s: %w", filename, err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON data to file %s: %w", filename, err)
	}

	fmt.Printf("Info: Successfully saved data to %s.\n", filename)
	return nil
}

func (repo *linkRepository) LoadDataFromJSON() error {
	var (
		hashToLinkNumMap  map[uint64]int
		linkNumToLinksMap map[int][]string
	)

	// Load mapHashToLinkNum
	if err := loadFileToMap("./data/hash_to_link_num.json", &hashToLinkNumMap); err != nil {
		return fmt.Errorf("failed to load hash_to_link_num data: %w", err)
	}

	// Load mapLinkNumToLinks
	if err := loadFileToMap("./data/link_num_to_links.json", &linkNumToLinksMap); err != nil {
		return fmt.Errorf("failed to load link_num_to_links data: %w", err)
	}

	// Load data from maps into sync.Map's
	for hash, linkNum := range hashToLinkNumMap {
		repo.mapHashToLinkNum.Store(hash, linkNum)
	}
	for linkNum, links := range linkNumToLinksMap {
		repo.mapLinkNumToLinks.Store(linkNum, links)
	}

	// find max linkNum to set the counter correctly
	maxLinkNum := 0
	for linkNum := range linkNumToLinksMap {
		if linkNum > maxLinkNum {
			maxLinkNum = linkNum
		}
	}
	repo.counter.Store(uint64(maxLinkNum) + 1)

	return nil
}

func loadFileToMap(filename string, target interface{}) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Info: Data file %s not found, skipping load.\n", filename)
		return nil
	}

	fileData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read data file %s: %w", filename, err)
	}

	if err = json.Unmarshal(fileData, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON data from %s: %w", filename, err)
	}

	fmt.Printf("Info: Successfully loaded data from %s.\n", filename)
	return nil
}
