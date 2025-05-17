package http

import (
	"go-template/internal/common"
	"net/http"
)

func postRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}
	err := common.App.ReadJSON(w, r, &input)
	if err != nil {
		common.App.ErrorResponse(w, r, 500, err.Error())
		return
	}

	id, err := common.App.Repository.CreateRepository(input.Name)
	if err != nil {
		common.App.ServerErrorResponse(w, r, err)
		return
	}

	err = common.App.WriteJSON(w, 201, common.Envelope{"id": id}, nil)
	if err != nil {
		common.App.ErrorResponse(w, r, 500, err.Error())
		return
	}
}

func getRepoByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := common.App.ReadIdParam(r, "repoId")
	if err != nil {
		common.App.NotFoundResponse(w, r)
		return
	}

	repository, err := common.App.Repository.GetRepositoryById(id)
	if err != nil {
		common.App.ServerErrorResponse(w, r, err)
		return
	}

	err = common.App.WriteJSON(w, 200, common.Envelope{"data": repository}, nil)
	if err != nil {
		common.App.ErrorResponse(w, r, 500, err.Error())
		return
	}
}
