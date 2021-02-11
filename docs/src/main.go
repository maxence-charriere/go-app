//go:generate go run gen/godoc.go
//go:generate go fmt

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/maxence-charriere/go-app/v8/pkg/cli"
	"github.com/maxence-charriere/go-app/v8/pkg/errors"
	"github.com/maxence-charriere/go-app/v8/pkg/logs"
)

type localOptions struct {
	Port int `cli:"p" env:"GOAPP_DOCS_PORT" help:"The port used by the server that serves the PWA."`
}

type githubOptions struct {
	Output string `cli:"o" env:"-" help:"The directory where static resources are saved."`
}

func main() {
	for path := range pages() {
		app.Route(path, newPage())
	}
	app.Route("/reference", newReference())
	app.RunWhenOnBrowser()

	ctx, cancel := cli.ContextWithSignals(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()
	defer exit()

	localOpts := localOptions{Port: 7777}
	cli.Register("local").
		Help(`Launches a server that serves the documentation app in a local environment.`).
		Options(&localOpts)

	githubOpts := githubOptions{}
	cli.Register("github").
		Help(`Generates the required resources to run the documentation app on GitHub Pages.`).
		Options(&githubOpts)

	backgroundColor := "#2e343a"

	h := app.Handler{
		Author:          "Maxence Charriere",
		BackgroundColor: backgroundColor,
		Description:     "Documentation for the go-app package.",
		Keywords: []string{
			"go-app",
			"go",
			"golang",
			"app",
			"pwa",
			"progressive web app",
			"webassembly",
			"web assembly",
			"gui",
			"ui",
			"user interface",
			"graphical user interface",
			"frontend",
			"opensource",
			"open source",
			"github",
		},
		LoadingLabel: "Loading go-app documentation...",
		Name:         "Go-app Docs",
		Scripts: []string{
			"/web/js/prism.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
			"https://fonts.googleapis.com/css2?family=Roboto&display=swap",
			"/web/css/prism.css",
			"/web/css/docs.css",
		},
		RawHeaders: []string{
			`
			<!-- Global site tag (gtag.js) - Google Analytics -->
			<script async src="https://www.googletagmanager.com/gtag/js?id=G-SW4FQEM9VM"></script>
			<script>
			  window.dataLayer = window.dataLayer || [];
			  function gtag(){dataLayer.push(arguments);}
			  gtag('js', new Date());
			
			  gtag('config', 'G-SW4FQEM9VM');
			</script>`,
			`<meta property="og:url" content="https://go-app.dev">`,
			`<meta property="og:title" content="go-app">`,
			`<meta property="og:description" content="A package to build progressive web apps with Go programming language and WebAssembly.">`,
			`<meta property="og:image" content="https://go-app.dev/web/images/go-app.png">`,
		},
		ThemeColor: backgroundColor,
		Title:      "go-app documentation",
	}

	switch cli.Load() {
	case "local":
		runLocal(ctx, &h, localOpts)

	case "github":
		generateGitHubPages(ctx, &h, githubOpts)
	}
}

func runLocal(ctx context.Context, h *app.Handler, opts localOptions) {
	app.Log("%s", logs.New("starting go-app documentation service").
		Tag("port", opts.Port).
		Tag("version", h.Version),
	)

	s := http.Server{
		Addr:    fmt.Sprintf(":%v", opts.Port),
		Handler: h,
	}

	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func generateGitHubPages(ctx context.Context, h *app.Handler, opts githubOptions) {
	pages := pages()
	p := make([]string, 0, len(pages))
	for path := range pages {
		p = append(p, path)
	}

	if err := app.GenerateStaticWebsite(opts.Output, h, p...); err != nil {
		panic(err)
	}
}

func exit() {
	err := recover()
	if err != nil {
		app.Log("command failed: %s", errors.Newf("%v", err))
		os.Exit(-1)
	}
}
