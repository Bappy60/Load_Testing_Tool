package services

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/Bappy60/BookStore_in_Go/pkg/domain"
	"github.com/Bappy60/BookStore_in_Go/pkg/models"
	"github.com/Bappy60/BookStore_in_Go/pkg/types"
	"github.com/redis/go-redis/v9"
)

type BookService struct {
	repo domain.IBookRepo
	redisClient *redis.Client
}

func BookServiceInstance(bookRepo domain.IBookRepo,redisClient *redis.Client) domain.IBookService {
	return &BookService{
		repo: bookRepo,
		redisClient: redisClient,
	}
}

var ctx = context.Background()

func (service *BookService) GetBooks(reqStruct *types.BookReqStruc) ([]models.Book, error) {


	parsedId, err := strconv.ParseInt(reqStruct.ID, 0, 0)
	if err != nil && reqStruct.ID != "" {
		return nil,err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil,err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil,err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil,err
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
		return nil,err
	}
	book := &models.Book{
		ID: uint(bookID),
		Name: reqBook.Name,
		NumberOfPages: reqBook.NumberOfPages,
		Publication: reqBook.Publication,
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
		return "invalid format of ID",err
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
		return nil,err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil,err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil,err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil,err
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

func (service *BookService) GetBooksFromMap(reqStruct *types.BookReqStruc) ([]models.Book, error) {


	parsedId, err := strconv.ParseInt(reqStruct.ID, 0, 0)
	if err != nil && reqStruct.ID != "" {
		return nil,err
	}
	parsedAuthorID, err := strconv.ParseUint(reqStruct.AuthorID, 0, 0)
	if err != nil && reqStruct.AuthorID != "" {
		return nil,err
	}
	parsedNumOfPages, err := strconv.ParseInt(reqStruct.NumberOfPages, 0, 0)
	if err != nil && reqStruct.NumberOfPages != "" {
		return nil,err
	}
	parsedPublicationYear, err := strconv.ParseInt(reqStruct.PublicationYear, 0, 0)
	if err != nil && reqStruct.PublicationYear != "" {
		return nil,err
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