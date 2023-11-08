package abstract

import "fmt"

// 抽象工厂方法

type Person interface {
	Greet()
}

type person struct {
	name string
	age  int
}

func (p person) Greet() {
	fmt.Printf("Hi! My name is %s,", p.name)
}

// NewPerson 抽象工厂模式：返回一个接口
func NewPerson(name string, age int) Person {
	return person{
		name: name,
		age:  age,
	}
}
