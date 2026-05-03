package handler

import (
	"net/http"
	"strconv"

	"example.com/internal/dto"

	"github.com/gin-gonic/gin"
)

// parseIDParam extracts and validates a uint ID from the ":id" URL parameter.
// Returns the parsed ID and true on success, or sends an error response and returns false.
func parseIDParam(c *gin.Context) (uint64, bool) {
	raw := c.Param("id")
	if raw == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "ID is required", ""))
		return 0, false
	}

	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_ID", "Invalid ID format", ""))
		return 0, false
	}
	return id, true
}
