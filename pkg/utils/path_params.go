package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ParseIDParam(r *http.Request, paramName string) (int64, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[paramName]
	if !ok || idStr == "" {
		return 0, fmt.Errorf("missing %s", paramName)
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", paramName)
	}

	return id, nil
}
