package main

func main() {
	m := make(map[string]string)

	m["hi"] = "hi"
	m["hello"] += "world"
	m["hi"] += " yb"

	for k, v := range m {
		println(k, v)
	}
}
