package main

func main() {
	a := App{}
	a.Port = ":8080"
	a.Initialize()
	a.Run()
}
