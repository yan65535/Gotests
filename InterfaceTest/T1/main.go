package main

import "fmt"

type teacher interface { //定义通用接口类型 授课和检查代码两个方法
	speak() string
	codeReview() string
}
type aTeacher struct { //定义a教师
	aFunc string
}
type bTeacher struct { //定义b教师
	bFunc string
}
type tStruct struct { //定义教师
	teacher teacher
}

func (a *aTeacher) Aaa() string {
	fmt.Println("A技术")
	return a.aFunc + "A独有"
}

func (a aTeacher) speak() string {
	a.Aaa()
	return a.aFunc + "speakA"
}
func (a aTeacher) codeReview() string {
	return a.aFunc + "codeReviewA"
}
func (b bTeacher) Bbb() string {
	fmt.Println("B技术")
	return b.bFunc + "B独有"
}
func (b bTeacher) speak() string {
	b.Bbb()
	return b.bFunc + "speakB"
}
func (b bTeacher) codeReview() string {
	return b.bFunc + "codeReviewB"
}

func (t *tStruct) Build() {
	t.teacher.speak()      //教师授课
	t.teacher.codeReview() //教师检查代码
	fmt.Println("授课完成")
}

func main() {
	t := tStruct{aTeacher{}} // 启用a教师

	fmt.Println(t.teacher.speak())      // a教师授课
	fmt.Println(t.teacher.codeReview()) // a教师检查代码
	t.teacher = bTeacher{}              // 切换为b教师
	fmt.Println(t.teacher.speak())      // b教师授课
	fmt.Println(t.teacher.codeReview()) //b教师检查代码
	t.Build()                           //课程结束
}

//上面例子中，可以观察到接口的实现是隐式的，也对应了官方对于基本接口实现的定义：
//方法集是接口方法集的超集，所以在 Go 中，实现一个接口不需要implements关键字显式的去指定要实现哪一个接口，只要是实现了一个接口的全部方法，那就是实现了该接口。
//有了实现之后，就可以初始化接口了，教师结构体内部声明了一个Crane类型的成员变量，
//可以保存所有实现了Teacher接口的值，由于是Teacher 类型的变量，
//所以能够访问到的方法只有Speak 和CodeReview，内部的其他方法例如Aaa()和Bbb()都无法访问。
