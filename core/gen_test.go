package core

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

type K interface {
	~int | ~int64
	error
}
type Num int

func (n Num) Error() string {

	//TODO implement me
	return "nil"
}

//func (n Num) Clac[K](k1 k2 K)  {
//
//}

func clac[k ~int | ~float64](k1, k2 k) k {
	//k1.Error()
	fmt.Println(k1 == k2)
	return k1 + k2
}

type Pointer interface {
	*int
}

func pClac[P Pointer](pointer P) P {
	//var p P
	//return p
	return pointer
}

func TestPClac(t *testing.T) {
	a := 1
	t.Log(pClac(&a), pClac(&a) == nil)
}

func TestName(t *testing.T) {
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	t.Log(clac(Num(3), Num(3)))

	t.Log(SumIntsOrFloats(map[int64]int64{
		1: 2,
	}, 1))
	t.Log(SumIntsOrFloats[string, int64](ints, "first"))
	t.Log(SumIntsOrFloats(ints, "first"))
	t.Log(SumIntsOrFloats(floats, "first"))
}

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V, key K) V {

	return m[key]
}

type Stringer interface {
	String() string
}
type Str int

func (s Str) String() string {
	return strconv.Itoa(int(s))
}

type Plusser interface {
	Plus(string) string
}
type Plus string

func (p Plus) Plus(s string) string {
	return string(p) + s
}

func ConcatTo[S Stringer, P Plusser](s []S, p P) (res string) {

	for _, v := range s {
		res += p.Plus(v.String())
	}
	return
}
func TestConcat(t *testing.T) {
	t.Log(ConcatTo([]Str{122, 123, 344, 55}, Plus("-")))
}

type EmbededParamter[T any] interface {
	~int | ~float64
	me() T
}
type C int

func (n Num) me() Num {
	return n
}
func Abs[E EmbededParamter[E]](e E) int {
	fmt.Println(e == 1)
	fmt.Println(e) // 输出 nil
	//fmt.Println(e == nil) //invalid operation: e == nil (mismatched types E and untyped nil)
	fmt.Println(interface{}(e))

	return int(e.me())
}

func TestAbs(t *testing.T) {
	t.Log(Abs(Num(1)))
}

type Lockable[T any] struct {
	a T
	sync.Mutex
}

func TestLock(t *testing.T) {
	lock := Lockable[int]{
		a: 1,
	}
	lock.Lock()
	defer lock.Unlock()
	//lock.TryLock()
	t.Log(lock.a)
}

func TestMap(t *testing.T) {

}

type Float interface {
	~float32 | ~float64
}

func NewtonSqrt[T Float](v T) T {

	var iterations int

	switch (interface{})(v).(type) { //识别不了。。。。

	case float32:

		iterations = 4

	case float64:

		iterations = 5

	default:
		//走到这里了。。
		panic(fmt.Sprintf("unexpected type %T", v))

	}

	// Code omitted.
	return T(iterations)
}

type MyFloat float32

func TestNewtonSqrt(t *testing.T) {
	var _ = NewtonSqrt(MyFloat(64))
}

type Aget[T any] struct {
	t *T
}

// 根据实际判断，如果a的t不等于nil再返回，如果是nil就返回一个T类型的nil（意思就是只声明）

func (a *Aget[T]) Approach() (r T) {

	if a.t != nil {
		return *a.t

	}
	//var r T //零值只能这么声明

	return r

}

type Person struct {
	name string
	age  int
}

func TestApproach(t *testing.T) {
	n := 1
	var a = Aget[int]{
		t: &n,
	}

	t.Log(a.Approach())
	a.t = nil
	t.Log(a.Approach())

	var m = Aget[Person]{
		t: &Person{age: 1, name: "wfy"},
	}
	t.Log(m.Approach())
	m.t = nil
	t.Log(m.Approach())
}

type Request[T any] struct {
	param T
}
type F func(t interface{})

func (f F) Call(request Request[CommonReq]) {
	fmt.Println("before call")
	f(request)
	fmt.Println("after call")
}

type CommonReq struct {
	UID int64
}

func paramF[T func() int](t T) int {
	fmt.Println("before paramF")
	return t()
}

func TestParamF(t *testing.T) {
	t.Log(paramF(func() int {
		return 4
	}))
}
func TestFunc(t *testing.T) {
	var (
		f = func(t any) {
			fmt.Println("接受参数 t ", t)
			p := interface{}(t.(Request[CommonReq]).param)
			if req, ok := p.(CommonReq); ok {
				fmt.Println("commonReq ", req.UID)
			}
		}

		request = Request[CommonReq]{
			param: CommonReq{
				UID: 12030027,
			},
		}
	)
	F(f).Call(request)
}
