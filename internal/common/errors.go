package common

import (
	"fmt"
	"net/http"
)

// generic method to log error
func (app *Application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.Logger.Error(err.Error(), "method", method, "URI", uri)
}

// generic return error response
func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, message any) {
	envelope := Envelope{"error": message}

	err := app.WriteJSON(w, statusCode, envelope, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

}

func (app *Application) EditConflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "unable to update the record due to an edit conflict, please try again"
	app.ErrorResponse(w, r, http.StatusConflict, message)
}

func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *Application) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}
