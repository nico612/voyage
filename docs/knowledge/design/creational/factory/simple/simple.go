package simple

import "fmt"

// 简单工厂

type Person struct {
	Name string
	Age  int
}

func (p Person) Greet() {
	fmt.Printf("Hi! My name is %s,", p.Name)
}

// NewPerson 简单工厂模式
func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		Age:  age,
	}
}
