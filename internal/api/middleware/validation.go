package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

type ValidationResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

// ValidationMiddleware creates a middleware for validating ReceiptPayload.
func ValidationMiddleware[T any](validate *validator.Validate) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request T
		// Bind JSON to the payload
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, ValidationResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request format",
				Errors: []ValidationError{
					{
						Field: "body",
						Tag:   "json",
						Value: err.Error(),
					},
				},
			})
			ctx.Abort()
			return
		}

		// Validate the request
		if err := validate.Struct(&request); err != nil {
			var validationErrors []ValidationError

			// Cast err to validator.ValidationErrors to get detailed validation errors
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrs {
					validationErrors = append(validationErrors, ValidationError{
						Field: e.Field(),
						Tag:   e.Tag(),
						Value: e.Param(),
					})
				}
			}

			ctx.JSON(http.StatusBadRequest, ValidationResponse{
				Status:  http.StatusBadRequest,
				Message: "Validation failed",
				Errors:  validationErrors,
			})
			ctx.Abort()
			return
		}

		ctx.Set("payload", request)
		ctx.Next()
	}
}
