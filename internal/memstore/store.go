package memstore

import (
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID                              uuid.UUID
	Author, Title, Summary, Content string
	Source                          *url.URL
	Tags                            []string
	CreatedAt, UpdatedAt, DeletedAt time.Time
}

type Store struct {
	lock sync.RWMutex
	news []*News
}

func New() *Store {
	return &Store{
		news: make([]*News, 0),
		lock: sync.RWMutex{},
	}
}

func (s *Store) Create(news *News) *News {
	createdNews := &News{
		ID:        news.ID,
		Author:    news.Author,
		Title:     news.Title,
		Summary:   news.Summary,
		Content:   news.Content,
		Source:    news.Source,
		Tags:      news.Tags,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.news = append(s.news, createdNews)
	return createdNews
}

func (s *Store) Get(id uuid.UUID) *News {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, news := range s.news {
		if news.ID == id && news.DeletedAt.IsZero() {
			return news
		}
	}
	return nil
}

func (s *Store) GetAll() []*News {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.news
}

func (s *Store) Delete(id uuid.UUID) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, news := range s.news {
		if news.ID == id {
			//Set zero to timestamp
			news.DeletedAt = time.Time{}
		}
	}
}
