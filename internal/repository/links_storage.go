package repository

import (
	"Link-Status-Service/internal/utils"
	"bytes"
	"context"
	"errors"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

type linkRepository struct {
	counter              *atomic.Uint64
	map_hash_to_LinkNum  sync.Map
	map_LinkNum_to_links sync.Map
}

func NewLinkRepository() *linkRepository {
	repo := linkRepository{counter: &atomic.Uint64{}}
	repo.counter.Add(1)
	return &repo
}

func (repo *linkRepository) GetLinkNum(ctx context.Context, links []string) (linkNum int, isNew bool, err error) {
	sortedLinks := utils.SortStrings(links)
	hash64 := getHashOfLinks(sortedLinks)
	linkNum, isNew = repo.hashToLinkNum_GetOrCreate(hash64)
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

func (repo *linkRepository) hashToLinkNum_GetOrCreate(hash uint64) (linkNum int, isNew bool) {
	newLinkNum := repo.counter.Load()
	actualLinkNum, alreadyExisted := repo.map_hash_to_LinkNum.LoadOrStore(hash, int(newLinkNum))

	isNew = !alreadyExisted
	if isNew {
		repo.counter.Add(1)
	}
	return actualLinkNum.(int), isNew
}

func (repo *linkRepository) StoreLinks(ctx context.Context, links []string, linkNum int) error {
	sortedLinks := utils.SortStrings(links)
	repo.map_LinkNum_to_links.Store(linkNum, sortedLinks)

	return nil
}

func (repo *linkRepository) GetLinksByLinkNum(ctx context.Context, linkNum int) ([]string, error) {
	links, exists := repo.map_LinkNum_to_links.Load(linkNum)
	if !exists {
		return []string{}, nil
	}

	_, ok := links.([]string)
	if !ok {
		return []string{}, errors.New("cannot convert links to []string")
	}
	return links.([]string), nil
}

func (repo *linkRepository) StoreDataToJSON() error {
	return nil
}

func (repo *linkRepository) LoadDataFromJSON() error {
	return nil
}
