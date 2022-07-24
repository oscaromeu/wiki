---
id: functions_3
title: Function Values and Closures
sidebar_label: Functions 3 - Function Values and Closures
sidebar_position: 28
---

## Functions as values

Functions can be assigned to a variable just like an `int` or a `string` value. Function values may be used as function arguments and return values. The type of the variable is the function's signature, including parameter types and return types, but no parameter names.

```go
func f1(s string) bool { 
    return len(s) > 0 
}

func f2(s string) bool { 
    return len(s) < 4 
}

// Declare a variable of type func(string) bool
var funcVar func(string) bool

func main() {
    // Assign functions to the variable
    funcVar = f1
    fmt.Println(funcVar("abcd"))
    funcVar = f2
    fmt.Println(funcVar("abcd"))
}
```

## Function literals or anonymous functions

A function literal or anonymous function, is basically a function without name. An anonymous function is created dynamically, much like a variable is. 

```go
package main

import "fmt"

var DoStuff func() = func() {
  // Do stuff
}

func main() {
  DoStuff()

  DoStuff = func() {
    fmt.Println("Doing stuff!")
  }
  DoStuff()

  DoStuff = func() {
    fmt.Println("Doing other stuff.")
  }
  DoStuff()
}
```

Anonymous functions can accept parameters, return data, and do pretty much anything else a normal function can do. You can even assign a regular function to a variable just like you do with an anonymous function.

```go
package main

import "fmt"

var DoStuff func() = func() {
	// Do stuff
}

func RegFunc() { fmt.Println("reg func") }

func main() {
	DoStuff()
	DoStuff = RegFunc
	DoStuff()
}
```

## Passing functions to functions

A function can accept a function as a parameter. Like before, the parameter type is the function's signature but without any parameter names.

In this example, the type of `f` is `func(string) bool`

```go
func f1(s string) bool { 
    return len(s) > 0 
}

func f2(s string) bool { 
    return len(s) < 4 
}

func funcAsParam(s string, f func(string) bool) bool {
    return f(s + "abcd")
}

func main() {
    fmt.Println(funcAsParam("abcd", f1))
}
```

## Immediately calling a function literal

It is possible to define and call a function literal in one single statement by appending a parameter list `()` to the function literal. 

```go
var result string = func() string {
    return "abcd"
}()
```

## Closures

A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables. 

For example, the `adder` function returns a closure. Each closure is bound to its own `sum` variable.

```go
package main

import "fmt"

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func main() {
	pos := adder()
	for i := 0; i < 10; i++ {
		fmt.Println(pos(i))
	}
}
```

A closure has a very useful property: It has access to the local variables of the outer function, even after the outer function has terminated. Those variables then behave like static variables to the closure; in other words, they keep their values between subsequent calls to the closure. 

<figcaption align = "left"><b>Example. <a href="https://go.dev/play/p/nKb1PWQyl36">Go playground</a></b></figcaption>

```go
package main

import "fmt"

func newClosure() func() {
	var a int
	return func() {
		fmt.Println(a)
		a++
	}
}

func main() {
	c := newClosure()
	c()
	c()
	c()
}
```

### Closures and loop variables

The closures access the loop variable `c`, each of the closures should receive a different value of `c`, since `c` changes with every loop iteration. However, when running the second loop, all closures stored in the slice print out the same value. What happened?

```go
package main

import "fmt"

func caveat() {
	s := "abcd"

	var funcs []func()

	for _, c := range s {
		// c := c
		funcs = append(funcs, func() {
			fmt.Print(string(c))
		})
	}

	for _, f := range funcs {
		f()
	}
}

func main() {
	caveat()
}
```

<figcaption align = "left"><b>(<a href="https://go.dev/play/p/pXV5jBBLy7C">Go playground</a></b>)</figcaption>

Let's have the closures print out the address of `c` rather than its value. 

```go
funcs = append(funcs, func() {
    fmt.Println(&c)
})
```

All closures refer to the same instance of `c`! What we need instead is that each closure receives the value of `c` that exists at the point when the closure is created. We can achieve this declaring a variable within a loop body this way a new instance of this variable is created on every loop iteration. So when we change the first `for` loop to this:

```go
for _, c := range s {
    loopBodyVar := c
    funcs = append(funcs, func() {
        fmt.Print(string(loopBodyVar))
    })
}
```

Then each closure gets a separate instance of `loopBodyVar`. We can take this one step further and use the name of variable `c` also for the variable in the loop body:

```go
for _, c := range s {
    c := c
    funcs = append(funcs, func() {
        fmt.Print(string(c))
    })
}
```

Due to the scope rules, the `c` declared within the loop shadows the `c` declared in the loop condition, so even though this looks a bit strange it is perfectly valid code.

:::info
Take care when using the variable of a loop condition to define closures within the loop. The variable is the same instance for all closures created within the loop, hence all closures will read the same value when executing after the loop. Instead, make a copy of the loop condition variable within the loop body and have the closures uses that variable instead.
:::


## Summary

+ Functions are first-cass objects in Go. They can be assigned to variables and passed to other functions as parameters.
+ An anonymous function or function literal can be invoked directly after declaration by appending the parameter list to the literal. 
+ A closure is a function value that references variables from outside its body.


## Exercises

:::tip Exercises

1. We know that a closure can reference the outer function's variable even after the outer function has terminated. But what happens if the outer function generates and returns two closures? Do they access the same outer variables, or does each of them get its own copy?. Copy and paste the code from below into your editor and name the file "closures.go". 

    ```go
    package main
    import "fmt"
    func newClosures() (func(), func() int) {
        a := 0
        // Your code here!
    }
    func main() {
        f1, f2 := newClosures()
        f1() // sets "a" to 5
        n := f2() // multiplies "a" by 7 - is f2's internal value of "a" 0 or 5 before the call? 
        fmt.Println(n)
    }
    ```

   Add code to `newClosure` so that it returns two closures. The first one is of type `func()`, the second one is of type `func() int`. Both closures should modify an integer variable defined in the outer function as follows:

   + The first closure shoud just set the outer variable to 5. It returns nothing.
   + The second closure should multiply the outer variable by 7 and return the value. 
    
   The `main()` function calls `newClosure` to create the new closures, and then calls both closures and prints out the result. 
     
   <details>
     <summary>Toggle Solution</summary>
   
   ```go
      package main

      import "fmt"
      
      func newClosures() (func(), func() int) {
      	a := 0      
      	return func() { a = 5 }, func() int { return a * 7 }
      }

      func main() {
      	f1, f2 := newClosures()
      	f1()      // sets "a" to 5
      	n := f2() // multiplies "a" by 7 - is f2's internal value of "a" 0 or 5 before the call?
      	fmt.Println(n)
      }
   ```  
   </details>

   [https://go.dev/play/p/t6_5Jf9wawq](https://go.dev/play/p/t6_5Jf9wawq)

2. The `defer` keyword allows to specify a function that is called whenever the current function ends. What if we could call one function at the beggining of the current function, and one at the end, with only one call? Like so:

  ```go
  func f() {
      trace("f")
      fmt.Println("Doing something")
  }
  ``` 

  And when calling function f() it would print:

  ```
  Entering f
  Doing something
  Leaving f
  ```

  Write a function `trace()` that receives a string - the name of the current function - and does the following:

   + Print `Entering <name>` where `<name>` is the string parameter that trace receives
   + Create and return a function that prints `Leaving <name>`

  Then call trace() via the defer keyword in such a way that trace() runs immediately, and returns its result to defer.

  ```go
  func trace(name string) func() {
      // TODO:
      // 1. Print "Entering <name>"
      // 2. return a func() that prints "Leaving <name>" 
  }
  func f() {
      defer // TODO: add trace() here so the defer receives the returned function
      fmt.Println("Doing something")
  }
  func main() {
      fmt.Println("Before f")
      f()
      fmt.Println("After f")
  }
  ```

  <details>
    <summary>Toggle Solution</summary>

    ```go
    package main
    
    import "fmt"
    
    func trace(name string) func() {
    	fmt.Println("Entering", name)
    	return func() {
    		fmt.Println("Leaving", name)
    	}
    }
    func f() {
    	defer trace("f")()
    	fmt.Println("Doing something")
    }
    func main() {
    	fmt.Println("Before f")
    	f()
    	fmt.Println("After f")
    }
    ```

  </details>

  [https://go.dev/play/p/_7hx44LIuIS](https://go.dev/play/p/_7hx44LIuIS)

:::