package core

type Module interface {
	Run()
	CleanUp()
}
