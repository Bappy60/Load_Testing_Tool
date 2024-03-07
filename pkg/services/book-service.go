package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Bappy60/BookStore_in_Go/pkg/domain"
	"github.com/Bappy60/BookStore_in_Go/pkg/models"
	"github.com/Bappy60/BookStore_in_Go/pkg/types"
	"github.com/redis/go-redis/v9"
)

type BookService struct {
	repo        domain.IBookRepo
	redisClient *redis.Client
}

func BookServiceInstance(bookRepo domain.IBookRepo, redisClient *redis.Client) domain.IBookService {
	return &BookService{
		repo:        bookRepo,
		redisClient: redisClient,
	}
}

var ctx = context.Background()

func (service *BookService) GetBooks(reqStruct *types.BookReqStruc) ([]models.Book, error) {

	parsedId, err := strconv.ParseInt(reqStruct.ID, 0, 0)
	if err != nil && reqStruct.ID != "" {
		return nil, err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil, err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil, err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil, err
	}

	pAuthorID := uint(parsedAuthorID)
	fbookstruc := types.FilterBookStruc{
		ID:              uint(parsedId),
		Name:            &reqStruct.Name,
		PublicationYear: int(parsedPublicationYear),
		NumberOfPages:   int(parsedNumOfPages),
		AuthorID:        &pAuthorID,
		Publication:     &reqStruct.Publication,
	}

	res, err := service.repo.GetBooks(&fbookstruc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *BookService) CreateBook(book *types.CreateBookStruc) (*models.Book, error) {

	res, err := service.repo.CreateBook(book)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *BookService) UpdateBook(reqBook *types.UpdateBookStruc) (*models.Book, error) {
	bookID, err := strconv.ParseUint(reqBook.ID, 0, 0)
	if err != nil {
		return nil, err
	}
	book := &models.Book{
		ID:              uint(bookID),
		Name:            reqBook.Name,
		NumberOfPages:   reqBook.NumberOfPages,
		Publication:     reqBook.Publication,
		PublicationYear: reqBook.PublicationYear,
	}
	res, err := service.repo.UpdateBook(book)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (service *BookService) DeleteBook(bookId string) (string, error) {

	parsedId, err := strconv.ParseInt(bookId, 0, 0)
	if err != nil {
		return "invalid format of ID", err
	}
	res, err := service.repo.DeleteBook(parsedId)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (service *BookService) GetBooksFromRedis(reqStruct *types.BookReqStruc) ([]models.Book, error) {

	parsedId, err := strconv.ParseInt(reqStruct.ID, 0, 0)
	if err != nil && reqStruct.ID != "" {
		return nil, err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil, err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil, err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil, err
	}

	pAuthorID := uint(parsedAuthorID)
	fbookstruc := types.FilterBookStruc{
		ID:              uint(parsedId),
		Name:            &reqStruct.Name,
		PublicationYear: int(parsedPublicationYear),
		NumberOfPages:   int(parsedNumOfPages),
		AuthorID:        &pAuthorID,
		Publication:     &reqStruct.Publication,
	}

	cacheKey := fmt.Sprintf("%d", fbookstruc.ID)
	if products, err := service.CheckInCache(cacheKey); err == nil {
		return products, nil
	}

	res, err := service.repo.GetBooks(&fbookstruc)
	if err != nil {
		return nil, err
	}

	if err := service.SetInCache(res, cacheKey); err != nil {
		return nil, err
	}

	return res, nil
}

// Define a map to cache books
var bookCacheMap sync.Map

// Function to populate the map with books from the database
func PopulateBookCacheMap(repo domain.IBookRepo) {
	fbook := &types.FilterBookStruc{}
	books, err := repo.GetBooks(fbook)
	if err != nil {
		fmt.Println("Error loading books into cache:", err)
		return
	}
	for _, book := range books {
		bookCacheMap.Store(book.ID, book)
	}
}

func (service *BookService) GetBooksFromMap(reqStruct *types.BookReqStruc) ([]models.Book, error) {

	// Check if the book ID is present in the request and if it's in the cache
	if reqStruct.ID != "" {
		parsedID, err := strconv.ParseUint(reqStruct.ID, 10, 64)
		if err != nil {
			return nil, err
		}
		if cachedBook, ok := bookCacheMap.Load(uint(parsedID)); ok {
			return []models.Book{*cachedBook.(*models.Book)}, nil
		}
	}

	parsedId, err := strconv.ParseInt(reqStruct.ID, 0, 0)
	if err != nil && reqStruct.ID != "" {
		return nil, err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil, err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil, err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil, err
	}

	pAuthorID := uint(parsedAuthorID)
	fbookstruc := types.FilterBookStruc{
		ID:              uint(parsedId),
		Name:            &reqStruct.Name,
		PublicationYear: int(parsedPublicationYear),
		NumberOfPages:   int(parsedNumOfPages),
		AuthorID:        &pAuthorID,
		Publication:     &reqStruct.Publication,
	}

	res, err := service.repo.GetBooks(&fbookstruc)
	if err != nil {
		return nil, err
	}

	// Update the cache with the newly retrieved books
	for _, book := range res {
		bookCacheMap.Store(book.ID, book)
	}

	return res, nil
}

func (BookService *BookService) CheckInCache(cacheKey string) ([]models.Book, error) {
	cachedData, err := BookService.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var books []models.Book
		err := json.Unmarshal([]byte(cachedData), &books)
		if err != nil {
			return nil, err
		}
		return books, nil
	}
	return nil, err
}

func (BookService *BookService) SetInCache(books []models.Book, cacheKey string) error {
	jsonData, err := json.Marshal(books)
	if err != nil {
		return err
	}
	if _, err := BookService.redisClient.Set(ctx, cacheKey, jsonData, time.Hour).Result(); err != nil {
		return err
	}
	return nil
}
