# nav
An app to demonstrate navigation between components.

![hello](https://github.com/murlokswarm/app/wiki/assets/nav.gif)

## Build on mac
![hello](https://github.com/murlokswarm/app/wiki/assets/nav-mac.png)

```bash
# In $GOPATH/src/github.com/murlokswarm/app/examples/hello/bin/nav-mac

# Build app
go build

## Launch app
./nav-mac
```


## Build on web
![hello](https://github.com/murlokswarm/app/wiki/assets/nav-web.png)

```bash
# In $GOPATH/src/github.com/murlokswarm/app/examples/nav/bin/nav-web

# Build server and client
goapp web build

# Launch server
./nav-web

# Launch client
open http://localhost:7042     # MacOS
explorer http://localhost:7042 # Windows
xdg-open http://localhost:7042 # Linux
```