package main


type Employee interface {
	IceCreamMaker
}

type IceCreamMaker interface {
	Hello()
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	j.name = "Jerry"
}


func main() {
	var ben = &Jerry{}
	var employee Employee = ben
	employee.Hello()
}
