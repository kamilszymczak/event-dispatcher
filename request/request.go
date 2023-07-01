package request

type Requestable[T any] interface {
	GetUrl() string
	GetData() T
	SetData(T) T
}

// https://www.reddit.com/r/golang/comments/z51a46/optional_function_parameters_and_generics_for/
type Request[T any] struct {
	Url string
	Data T
}

func (r *Request[T]) GetUrl() string {
	return r.Url
}

func (r *Request[T]) GetData() T {
	return r.Data
}

func (r *Request[T]) SetData(data T) T {
	r.Data = data
	return data
}

func New[T any](url string) Requestable[T] {
	request := &Request[T]{
		Url: url,
	}
	return request
}