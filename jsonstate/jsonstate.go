//提供返回json对象的常用结构
package jsonstate

type BoolString struct{
	State	bool
	Data	string
}
type BoolInt	struct{
	State	bool
	Data	int
}
type BoolState	struct{
	State	bool
	Data	interface{}
}
type IntString	struct{
	State	int
	Data	string
}
type	IntBool	struct{
	State	int
	Data	string
}
type	IntState struct{
	State   int
	Data	interface{}
}
type	StringInt	struct{
	State   string
	Data   int
}
type    StringBool	struct{
	State   string
	Data   bool
}
type    StringState	struct{
	State   string
	Data   interface{}
}
