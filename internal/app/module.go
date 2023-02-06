package app

type Module interface {
	Run()
	CleanUp()
}
