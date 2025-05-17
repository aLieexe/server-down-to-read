package models

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type Book struct {
	Filename     string      `db:"filename" json:"filename"`
	Link         string      `db:"link" json:"link"`
	RepositoryId pgtype.UUID `db:"repository_id" json:"repository_id"`
	Id           pgtype.UUID `db:"id" json:"id"`
	LinkExpiry   time.Time   `db:"link_expiry" json:"link_expiry"`
}

type CreatedBook struct {
	Filename string    `json:"filename" db:"filename"`
	Link     string    `json:"link" db:"link"`
	BookId   uuid.UUID `json:"book_id" db:"book_id"`
}

type BookModel struct {
	Pool       *pgxpool.Pool
	S3         *minio.Client
	BucketName string
}

func (b *BookModel) GetBookById(id uuid.UUID) (Book, error) {
	query := `SELECT * FROM books where id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.Pool.Query(ctx, query, args)
	if err != nil {
		return Book{}, ErrQueryError
	}
	book, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Book])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return Book{}, ErrRecordNotFound
		default:
			return Book{}, err
		}
	}

	return book, nil
}

func (b *BookModel) updateLink(id uuid.UUID) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", id))
	bucketName := b.BucketName

	presignedUrl, err := b.S3.PresignedGetObject(context.Background(), bucketName, id.String(), time.Hour*24*7, reqParams)
	if err != nil {
		return presignedUrl, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `UPDATE books SET link = @link WHERE id = @id`
	args := pgx.NamedArgs{
		"link": presignedUrl,
		"id":   id,
	}

	b.Pool.Exec(ctx, query, args)
	return presignedUrl, nil
}

func (b *BookModel) AddBook(fileHeader *multipart.FileHeader, repoId uuid.UUID, file io.Reader) (uuid.UUID, error) {

	exists := b.CheckIfRepositoryExist(repoId)

	if !exists {
		return uuid.Nil, ErrRecordNotFound
	}
	id := uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bucketName := b.BucketName
	_, err := b.S3.PutObject(ctx, bucketName, id.String(), file, int64(fileHeader.Size), minio.PutObjectOptions{
		ContentType: "application/pdf",
	})

	if err != nil {
		return uuid.Nil, err
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileHeader.Filename))

	presignedUrl, err := b.S3.PresignedGetObject(ctx, bucketName, id.String(), time.Hour*24*7, reqParams)
	if err != nil {
		return uuid.Nil, err
	}

	query := `INSERT INTO books (filename, link, repository_id, id) VALUES (@filename, @link, @repository_id, @id)`
	args := pgx.NamedArgs{
		"filename":      fileHeader.Filename,
		"link":          presignedUrl,
		"repository_id": repoId,
		"id":            id.String(),
	}

	commandTag, err := b.Pool.Exec(context.Background(), query, args)
	if err != nil {
		return uuid.Nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return uuid.Nil, ErrRecordNotFound
	}

	return id, nil
}

func (b *BookModel) GetBooksByRepoId(repoId uuid.UUID) ([]CreatedBook, error) {
	query := `SELECT * FROM books where repository_id = @repository_id`
	args := pgx.NamedArgs{
		"repository_id": repoId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.Pool.Query(ctx, query, args)
	if err != nil {
		return []CreatedBook{}, err
	}

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return []CreatedBook{}, err
	}

	createdBooks := make([]CreatedBook, 0)

	for _, book := range books {

		parsedBookId, err := uuid.Parse(book.Id.String())

		//if not expired yet then no need to do anything
		if book.LinkExpiry.Before(time.Now()) {
			createdBooks = append(createdBooks, CreatedBook{
				Filename: book.Filename,
				Link:     book.Link,
				BookId:   parsedBookId,
			})
			continue
		}

		if err != nil {
			return []CreatedBook{}, err
		}

		url, err := b.updateLink(parsedBookId)
		if err != nil {
			return []CreatedBook{}, err
		}

		book.Link = url.String()
		createdBooks = append(createdBooks, CreatedBook{
			Filename: book.Filename,
			Link:     book.Link,
			BookId:   parsedBookId,
		})
	}

	return createdBooks, err

}

func (b *BookModel) CheckIfRepositoryExist(repoId uuid.UUID) bool {
	query := `SELECT id FROM repository where id = @id`
	args := pgx.NamedArgs{
		"id": repoId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := b.Pool.Query(ctx, query, args)
	if err != nil {
		return false
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return false
		default:
			return false
		}
	}

	return true

}
