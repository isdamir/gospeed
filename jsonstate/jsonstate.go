//提供返回json对象的常用结构
package jsonstate

type BoolString struct{
	S	bool
	T	string
}
type BoolInt	struct{
	S	bool
	T	int
}
type BoolState	struct{
	S	bool
	T	interface{}
}
type IntString	struct{
	S	int
	T	string
}
type	IntBool	struct{
	S	int
	T	string
}
type	IntState struct{
	S   int
	T	interface{}
}
type	StringInt	struct{
	S   string
	T   int
}
type    StringBool	struct{
	S   string
	T   bool
}
type    StringState	struct{
	S   string
	T   interface{}
}
