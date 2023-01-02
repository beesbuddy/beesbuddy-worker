package core

type CmdRunner interface {
	Run()
	CleanUp()
}
