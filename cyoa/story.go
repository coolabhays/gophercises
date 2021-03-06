package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultStoryTmpl))
}

var tpl *template.Template

var defaultStoryTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<section class="page">
			<h1>{{.Title}}</h1>
			{{range .Paragraphs}}
			<p>{{.}}</p>
			{{end}}
			<br>
			<ul>
				{{range .Options}}
				<li><a href="/{{.Arc}}">{{.Text}}</a></li>
				{{end}}
			</ul>
		</section>
		<style>
			body {
				font-family: Source Code Pro;
			}
			h1 {
				text-align: center;
				position: relative;
			}
			.page {
				width: 80%;
				max-width: 750px;
				margin: auto;
				margin-top: 40px;
				margin-bottom: 40px;
				padding: 40px;
				background: #e2ddff;
				border: 1px solid #eee;
				box-shadow: 5px 10px #333;
			}
			ul {
				border-top: 1px dotted #ccc;
				padding: 10px 0 0 0;
				-webkit-padding-start: 0;
			}
			li {
				padding: 10px;
			}
			a,
			a:visited {
				text-decoration: none;
				color: #659ef8;
			}
			a,
			a:hover {
				color: #319aff;
			}
			p {
				text-indent: 1em;
			}
		</style>
	</body>
</html>`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFunc(fn func(req *http.Request) string) HandlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

// functional options
func NewHandler(story Story, opts ...HandlerOption) http.Handler {
	h := handler{story, tpl, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
	pathFn func(req *http.Request) string
}

func defaultPathFn(req *http.Request) string {
	path := strings.TrimSpace(req.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := h.pathFn(req)

	if chapter, ok := h.s[path]; ok {
		if err := h.t.Execute(w, chapter); err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "chapter not found.", http.StatusNotFound)
}

type Chapter struct {
	Title string `json:"title"`
	Paragraphs []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc string `json:"arc"`
}

type Story map[string]Chapter

func ReadStory(data io.Reader) (Story, error) {
	decoder := json.NewDecoder(data)
	var story Story
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
