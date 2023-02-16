package internal

type Ctx interface {
	Run()
	Flush()
}
