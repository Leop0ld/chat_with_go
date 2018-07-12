package trace

import (
	"fmt"
	"io"
)

// Tracer 는 코드 전체에서 이벤트를 추적할 수 있는 객체를 설명하는 인터페이스
type Tracer interface {
	Trace(...interface{})
}

// New 함수는 새로운 Tracer 생성함
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}
