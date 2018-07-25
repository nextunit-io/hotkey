# Hotkey

This repository contains a module to register hotkeys in your OS.

## Supported OS

We want provide some more OS. At the moment there are only the following available and the crossed are in doing.

* Windows
* ~~MAC OS~~
* ~~Linux~~

## Getting Started

To use hotkey it is very ease. Just import the module

```Go
import "github.com/nextunit-io/hotkey"
```

and in your code you can register your hotkeys. This provides an example for the `main` function.

```Go
package main

import (
	"fmt"

	"github.com/nextunit-io/hotkey"
)

func main() {
	defer hotkey.Close()

	h1, err1 := hotkey.Create(1, hotkey.ModAlt+hotkey.ModCtrl, 'O')
	h2, err2 := hotkey.Create(2, hotkey.ModAlt+hotkey.ModCtrl, 'X')

	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}

	endlessLoop := true

	h1.Register(func(id int) {
		fmt.Println("Hotkey pressed: ", h1.String())
	})
	h2.Register(func(id int) {
		fmt.Println("Hotkey pressed: ", h1.String())
		endlessLoop = false
	})

	for endlessLoop {
        // Here is some code of you - or whatever you wanna do :-)
	}
}
```