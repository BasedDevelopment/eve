package routes

import (
	"fmt"
	"net/http"

	"github.com/ericzty/eve/internal/controllers"
)

func GetHVs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(&controllers.Cloud)
}
