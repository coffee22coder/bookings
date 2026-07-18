package port

type GeneratorID interface {
	BookRef() (string, error)
}
