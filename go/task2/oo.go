package main

import "fmt"

type Shape interface {
	Area() float32
	Perimeter() float32
}

type Rectangle struct {
	side1 float32
	side2 float32
}

func (p *Rectangle) Area() float32 {
	return p.side1 * p.side2
}

func (p *Rectangle) Perimeter() float32 {
	return 2 * (p.side1 + p.side2)
}

type Circle struct {
	r float32
}

func (p *Circle) Area() float32 {
	return float32(3.14) * p.r * p.r
}

func (p *Circle) Perimeter() float32 {
	return 2 * float32(3.14) * p.r
}

func main() {
	rectangle := &Rectangle{float32(3), float32(4)}
	rectangleArea := rectangle.Area()
	rectanglePerimeter := rectangle.Perimeter()
	fmt.Printf("rectangle area: %v, perimeter: %v\n", rectangleArea, rectanglePerimeter)

	circle := &Circle{float32(3)}
	circleArea := circle.Area()
	circlePerimeter := circle.Perimeter()
	fmt.Printf("circle area: %v, perimeter: %v\n", circleArea, circlePerimeter)

	// ==================================
	e := &Employee{Person: Person{Name: "Bob", Age: 32}, EmployeeID: "9797"}
	e.PrintInfo()
}

type Person struct {
	Name string
	Age  uint8
}

type Employee struct {
	Person
	EmployeeID string
}

func (e *Employee) PrintInfo() {
	fmt.Printf("EmployeeID: %s, Name: %s, Age: %v", e.EmployeeID, e.Name, e.Age)
}
