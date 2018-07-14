
2.6.1 / 2018-07-13
==================

* interfaces and funcs renamed with shorter names:
  * interfaces:
    * `Component` => `Compo`
    * `ElementWithComponent` => `ElemWithCompo`
    * `ComponentWithExtendedRender` => `CompoWithExtendedRender`
  * functions:
    * `ComponentNameFromURLString` => `CompoNameFromURLString`

2.6.0 / 2018-07-12
==================

* `app.Element` is renamed `app.Elem`

* func to get an element changed.

  Old way (deprecated/deleted):

  ```go
  func ElementByCompo(c Compo) (ElemWithCompo, error)        { ... }
  func NavigatorByCompo(c Compo) (Navigator, error)          { ... }
  func WindowByCompo(c Compo) (Window, error)                { ... }
  func PageByCompo(c Compo) (Page, error)                    { ... }
  ```

  New way:

  ```go
  func ElemByCompo(c Compo) Elem { ... }
  ```

  Example:

  ```go
  app.ElemByCompo(myCompo).WhenWindow(func(w app.Window) {
    w.Center()
  })
  ```

  See [Elem](https://godoc.org/github.com/murlokswarm/app#Elem) definition for more info.

* Add compatibility with 10.11 (#144)

* MacOS driver is now compatible with MacOS El Capitan (10.11). Thanks to [tectiv3](https://github.com/tectiv3)

2.5.1 / 2018-06-07
==================

* Fix mac dock (#140)

2.5.0 / 2018-06-05
==================

* Status menu (#139)

2.4.3 / 2018-05-27
==================

* fix + travis (#137)

2.4.2 / 2018-05-26
==================

* Refactor logs decorators (#135)

2.4.1 / 2018-05-15
==================

* goapp mac

2.4.0 / 2018-03-25
==================

* Event and subscriber
* Actions

2.3.3 / 2018-03-23
==================

* Refactor decorators func

2.3.2 / 2018-03-18
==================

* bridge go rpc refactor

2.3.1 / 2018-03-16
==================

* bridge platform rpc refactor

2.3.0 / 2018-03-09
==================

* web driver implementation
* goapp tool to build web app

2.2.0 / 2018-02-25
==================

* Drag and drop

2.1.0 / 2018-02-24
==================

* Save file panel
* InnerText in HTML event handler when contenteditable is set

2.0.0 / 2018-02-23
==================

* Use of standard HTML in templates
* Get rid of hidden imports
* Improve code quality
* Improve architecture to make multiplaform easier to implement
* Centralize app related code into a single package (drivers, markup)

1.0.1 / 2018-02-03
==================
  
* Add history.md

1.0.0/ 2018-02-03
==================
  
* V1 release
