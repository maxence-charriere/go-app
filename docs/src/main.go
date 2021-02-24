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

const (
	defaultTitle       = "A Golang package for building progressive web apps"
	defaultDescription = "A package for building progressive web apps (PWA) with the  Go programming language and WebAssembly. It uses a declarative syntax that allows creating and dealing with HTML elements only by using Go, and without writing any HTML markup."
	backgroundColor    = "#2e343a"

	buyMeACoffeeURL   = "https://www.buymeacoffee.com/maxence"
	openCollectiveURL = "https://opencollective.com/go-app"
	githubURL         = "https://github.com/maxence-charriere/go-app"
	githubSponsorURL  = "https://github.com/sponsors/maxence-charriere"
	twitterURL        = "https://twitter.com/jonhymaxoo"
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
	app.Route("/issue499", newIssue499Data())
	app.Route("/", newMarkdownDoc())
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

	h := app.Handler{
		Name:        "Go-app Docs",
		Title:       defaultTitle,
		Description: defaultDescription,
		Author:      "Maxence Charriere",
		Image:       "/web/images/go-app.png",
		Keywords: []string{
			"go-app",
			"go",
			"golang",
			"app",
			"pwa",
			"progressive web app",
			"webassembly",
			"web assembly",
			"webapp",
			"web",
			"gui",
			"ui",
			"user interface",
			"graphical user interface",
			"frontend",
			"opensource",
			"open source",
			"github",
		},
		BackgroundColor: backgroundColor,
		ThemeColor:      backgroundColor,
		LoadingLabel:    "Loading go-app documentation...",
		Scripts: []string{
			"/web/js/prism.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
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
		},
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
