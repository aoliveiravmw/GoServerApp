package main

func main() {
	a := App{}
	a.Port = ":8081"
	a.Initialize()
	a.Run()
}
