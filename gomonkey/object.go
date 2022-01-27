package main

type privateInfo interface {
	GetName() string
	GetAge() int
}

type Info interface {
	GetName() string
	GetAge() int
}

type privateObject struct {
	name string
	age  int
}

func (p *privateObject) GetName() string {
	return p.name
}

func (p privateObject) GetAge() int {
	return p.age
}

func (p *privateObject) setName(name string) {
	p.name = name
}

type privateObjectWithPublicFields struct {
	Name string
	Age  int
}

func (p *privateObjectWithPublicFields) GetName() string {
	return p.Name
}

func (p privateObjectWithPublicFields) GetAge() int {
	return p.Age
}

type Object struct {
	Name string
	Age  int
}

func (p *Object) GetName() string {
	return p.Name
}

func (p Object) GetAge() int {
	return p.Age
}

type ObjectWithPrivateFields struct {
	name string
	age  int
}

func (p *ObjectWithPrivateFields) GetName() string {
	return p.name
}

func (p ObjectWithPrivateFields) GetAge() int {
	return p.age
}
