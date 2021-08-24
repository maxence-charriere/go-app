package main

type status int

const (
	neverLoaded status = iota
	loading
	loadingErr
	loaded
)
