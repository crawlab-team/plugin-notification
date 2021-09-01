package main

func main() {
	svc := NewService()
	if err := svc.Start(); err != nil {
		panic(err)
	}
}
