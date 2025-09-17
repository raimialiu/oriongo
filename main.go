package main

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	err := NewProgram().StartApp().Run()
	if err != nil {
		panic(err)
	}
}

// gorm(db), concurrency, advanced, (fx, log, test, etc, others) .....
