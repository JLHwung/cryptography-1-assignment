package week3

import (
	"fmt"
)

func ExampleFileAuth() {
	fmt.Printf("%x", FileAuth("./6.2.birthday.mp4_download"))
	// output: 03c08f4ee0b576fe319338139c045c89c3e8e9409633bea29442e21425006ea8
}
