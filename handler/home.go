package handler

import (
	"io/ioutil"
	"net/http"
)

//HomePageHandler 主页
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/view/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
