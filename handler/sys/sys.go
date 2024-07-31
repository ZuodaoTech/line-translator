package sys

import (
	"net/http"

	"github.com/zuodaotech/line-translator/handler/render"
)

func RenderRobots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// disable indexing
		render.Text(w, []byte("User-agent: *\nDisallow: /\n"))
	}
}

func RenderRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Html(w, []byte(`
			<html><head><title>Line Translator</title></head>
				<body>
					<p>Line Translator service, standing by.</p>
				</body>
			</html>`))
	}
}
