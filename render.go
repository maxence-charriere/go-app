package app

import (
	"github.com/murlokswarm/log"
	"github.com/murlokswarm/markup"
)

var (
	renderC     = make(chan markup.Componer)
	renderStopC = make(chan bool)
)

func startPipeRendering() {
	go func() {
		for {
			select {
			case c := <-renderC:
				render(c)

			case <-renderStopC:
				return
			}
		}
	}()
}

func stopPipeRendering() {
	renderStopC <- true
}

func render(c markup.Componer) {
	elems, err := markup.Sync(c)
	if err != nil {
		log.Panic(err)
	}

	for _, elem := range elems {
		driver.Render(elem.ID, elem.HTML())
	}
}
