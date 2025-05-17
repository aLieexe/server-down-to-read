package http

import (
	"errors"
	"go-template/internal/common"
	"go-template/internal/models"
	"net/http"
)

func postBookToRepoHandler(w http.ResponseWriter, r *http.Request) {
	//limit only to 16mb
	// multiform version
	r.Body = http.MaxBytesReader(w, r.Body, 16<<20)
	if err := r.ParseMultipartForm(16 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	//read Id
	repoId, err := common.App.ReadIdParam(r, "repoId")
	if err != nil {
		common.App.NotFoundResponse(w, r)
		return
	}

	bookId, err := common.App.Books.AddBook(fileHeader, repoId, file)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			common.App.NotFoundResponse(w, r)
			return
		}
		common.App.ServerErrorResponse(w, r, err)
		return
	}

	err = common.App.WriteJSON(w, 201, common.Envelope{"data": fileHeader.Filename, "bookId": bookId}, nil)
	if err != nil {
		common.App.ServerErrorResponse(w, r, err)
	}

}

func getBooksByRepoIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := common.App.ReadIdParam(r, "repoId")
	if err != nil {
		common.App.NotFoundResponse(w, r)
		return
	}

	books, err := common.App.Books.GetBooksByRepoId(id)
	if err != nil {
		common.App.ServerErrorResponse(w, r, err)
		return
	}

	err = common.App.WriteJSON(w, 200, common.Envelope{"books": books}, nil)
	if err != nil {
		common.App.ServerErrorResponse(w, r, err)
		return
	}
}
