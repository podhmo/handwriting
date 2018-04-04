package main

import (
	"fmt"
	"log"

	"github.com/podhmo/handwriting"
)

func main() {
	p := handwriting.Must(handwriting.New("./main"))
	f := p.File("main.go")
	f.Import("fmt")
	f.Code(func(f *handwriting.File) error {
		o := f.Out
		o.WithBlock("func main()", func() {
			o.WithBlock("for i := 1; i<=100; i++", func() {
				o.WithBlock("switch i", func() {
					for i := 0; i <= 100; i++ {
						o.WithIndent(fmt.Sprintf("case %d:", i), func() {
							if i%3 == 0 && i%5 == 0 {
								o.Println(`fmt.Println("fizzbuzz")`)
							} else if i%3 == 0 {
								o.Println(`fmt.Println("fizz")`)
							} else if i%5 == 0 {
								o.Println(`fmt.Println("buzz")`)
							} else {
								o.Println(`fmt.Printf("%d\n", i)`)
							}
						})
					}
					o.WithIndent("default:", func() {
						o.Println(`panic("not supported")`)
					})
				})
			})
		})
		return nil
	})
	if err := p.Emit(); err != nil {
		log.Fatal(err)
	}
}
