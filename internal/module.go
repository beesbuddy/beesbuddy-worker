package internal

type ModuleCtx interface {
	Run()
	CleanUp()
}
