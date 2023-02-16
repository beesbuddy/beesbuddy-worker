package internal

type Ctx interface {
	Run()
	CleanUp()
}
