package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)
import "github.com/starfederation/datastar-go/datastar"

import "color_picker/cmd/types"
import "color_picker/ui/templates/homepage"

func main() {
	static_fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", static_fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/homepage", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/homepage", get_home_page)
	http.HandleFunc("/get_rando_color", get_rando_color)
	http.HandleFunc("/roulette", get_roulette)

	fmt.Println("Server booting up on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func get_home_page(w http.ResponseWriter, r *http.Request) {
	rand_color_container := get_random_color_container()
	component := homepage_templ.Setup(rand_color_container)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}

func get_rando_color(w http.ResponseWriter, r *http.Request) {
	rando_container := get_random_color_container()
	component := homepage_templ.Get_clr_container(rando_container)

	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(component)
}

func get_random_color_container() types.Color_container {
	var res types.Color_container
	var cl types.Color

	for i := 0; i < 5; i++ {
		cl.R = uint8(rand.Uint32() % 256)
		cl.G = uint8(rand.Uint32() % 256)
		cl.B = uint8(rand.Uint32() % 256)

		res.Colors[i] = cl
	}
	return res
}

var frame_times = []float32{0, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14} // this is in ns

func get_roulette(w http.ResponseWriter, r *http.Request) {
	s_ind := r.FormValue("ind")
	ind, _ := strconv.Atoi(s_ind)
	ind++

	var frame_time float32 = 0

	if ind < len(frame_times) {
		frame_time = frame_times[ind]
		fmt.Printf("index: %d\n", ind)
	} else {
		fmt.Printf("len: %d, index: %d\n", len(frame_times), ind)
		ind = 0
	}

	rando_container := get_random_color_container()
	component := homepage_templ.Roulette(ind, frame_time, rando_container)

	time.Sleep(time.Millisecond * time.Duration(frame_time))
	sse := datastar.NewSSE(w, r)
	sse.PatchElementTempl(component)
}
