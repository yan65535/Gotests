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

func (a aTeacher) speak() string {
	return a.aFunc + "speakA"
}
func (a aTeacher) codeReview() string {
	return a.aFunc + "codeReviewA"
}

func (b bTeacher) speak() string {
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
	t := &tStruct{aTeacher{}}           // 启用a教师
	fmt.Println(t.teacher.speak())      // a教师授课
	fmt.Println(t.teacher.codeReview()) // a教师检查代码
	t = &tStruct{bTeacher{}}            // 切换为b教师
	fmt.Println(t.teacher.speak())      // b教师授课
	fmt.Println(t.teacher.codeReview()) //b教师检查代码
	t.Build()                           //课程结束
}
