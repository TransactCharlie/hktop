package widgets

type K8SWidget interface {
	Update() error
	Run()
	Stop() bool
}