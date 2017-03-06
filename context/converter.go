package context

type Converter interface {
	Convert(in interface{}) interface{}
}
