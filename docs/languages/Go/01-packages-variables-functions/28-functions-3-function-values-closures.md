---
id: functions_3
title: Function Values and Closures
sidebar_label: Functions 3 - Function Values and Closures
sidebar_position: 28
---

## Functions as values

Functions are values too and can be assigned to a variable just like an `int` or a `string` value. Function values may be used as function arguments and return values. The type of the variable is the function's signature, including parameter types and return types, but no parameter names.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/nDrnu3Ee7gD">https://go.dev/play/p/nDrnu3Ee7gD</a></b></figcaption>

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

A function literal or anonymous function, is basically a function without name. Anonymous functions can accept parameters, return data, and do pretty much anything else a normal function can do. 

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/WPHkRpCzH4f">https://go.dev/play/p/WPHkRpCzH4f</a></b></figcaption>


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

Notice that we create a variable that has a `func()` type. After doing this, we could create and assign an anonymous function to the variable. The big different with a regular function is that we could assign a new function to the `DoStuff` variable at runtime, allowing us to dynamically change what `DoStuff()` does. 


### Immediately calling a function literal

It is possible to define and call a function literal in one single statement by appending a parameter list `()` to the function literal. 

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/QtRJuf7n4G2">https://go.dev/play/p/QtRJuf7n4G2</a></b></figcaption>

```go
func main() {
	func() {
		fmt.Println("Greeting")
	}()
}
```

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/uEPg_pDdCnz">https://go.dev/play/p/uEPg_pDdCnz</a></b></figcaption>

```go
func main() {
	message := "Greeting"
	func(str string) {
		fmt.Println(str)
	}(message)
```

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/K1mNt6iUd7I">https://go.dev/play/p/K1mNt6iUd7I</a></b></figcaption>

```go
var result string = func() string {
    return "abcd"
}()
```

## Passing functions to functions

A function can accept a function as a parameter. Like before, the parameter type is the function's signature but without any parameter names.

In this example, the type of `f` is `func(string) bool`

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/BHMH-5PDIps">https://go.dev/play/p/BHMH-5PDIps</a></b></figcaption>

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



## Closures

A closure is a function value that references variables from outside its body. The function may access and assign to the referenced variables; in this sense the function is "bound" to the variables. 

For example, the `adder` function returns a closure. Each closure is bound to its own `sum` variable.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/QomfqGU_WKH">https://go.dev/play/p/QomfqGU_WKH</a></b></figcaption>

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



<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/nKb1PWQyl36">https://go.dev/play/p/nKb1PWQyl36</a></b></figcaption>

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

When you declare a new variable inside a `for` loop it is important to remember that the variables aren't being redeclared with each iteration. Instead the variable is the same, but instead the value that is stored in the variable is being updated. 

Let's look at the following example

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/pXV5jBBLy7C">https://go.dev/play/p/pXV5jBBLy7C</a></b></figcaption>

```go
package main

import "fmt"

func main() {
  var functions []func()

  for i := 0; i < 10; i++ {
    functions = append(functions, func() {
      fmt.Println(i)
    })
  }

  for _, f := range functions {
    f()
  }
}
```

The output we will get is:

```
10
10
10
10
10
10
10
10
10
10
```

The issue we are experiencing here is that `i` is declared inside a `for` loop, and it is being changed with each iteration of the for loop. When we finally call all of our functions in the `functions`slide they are all referencing the same `i` variable which was set to `10` in the last iteration of the for loop. The same thing happen if we use ranges instead.

So, how do we fix it? One way is to utilize the fact that function parameters in Go are passed by value. 

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/Of8UU_sQX-w">https://go.dev/play/p/Of8UU_sQX-w</a></b></figcaption>

```go
package main

import "fmt"

func main() {
  var functions []func()

  for i := 0; i < 10; i++ {
    functions = append(functions, build(i))
  }

  for _, f := range functions {
    f()
  }
}

func build(val int) func() {
  return func() {
    fmt.Println(val)
  }
}
```

Unfortunately, this example required us to create the `build()` function globally. Luckily there are some other solutions to solve this issue.

Here is the same example, but using an anonymous function to create the closure.

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/02KX3KIngds">https://go.dev/play/p/02KX3KIngds</a></b></figcaption>

```go
package main

import "fmt"

func main() {
  var functions []func()

  for i := 0; i < 10; i++ {
    functions = append(functions, func(val int) func() {
      return func() {
        fmt.Println(val)
      }
    }(i))
  }

  for _, f := range functions {
    f()
  }
}

```

+ First we declare a function inline that takes an integer value and returns a function.

  ```go
  func(val int) func() {
    return func() {
      fmt.Println(val)
    }
  }
  ```

+ Call the function with `i` as the parameter being passed in. This is the `(i)` part right after the function declaration.

+ After the anonymous function is called it returns a `func()`, which is then appended to the `functions` slice with the line `function=append(functions, ...)`

The last example can be written as

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/-obVXgm_nAq">https://go.dev/play/p/-obVXgm_nAq</a></b></figcaption>

```go
package main

import "fmt"

func main() {
  var functions []func()
  fn := func(val int) func() {
    return func() {
      fmt.Println(val)
    }
  }

  for i := 0; i < 10; i++ {
    functions = append(functions, fn(i))
  }

  for _, f := range functions {
    f()
  }
}
```

Notice how we are not adding `fn` to the `functions` slice, but we are passing in the return value of it which is a `func()`. Finally, we can also solve this problem by creating a new variable and assigning it with the value of `i`. Below is an example of this approach

<figcaption align = "left"><b>Example <a href="https://go.dev/play/p/Zgv2NCP6DCY">https://go.dev/play/p/Zgv2NCP6DCY</a></b></figcaption>

```go
package main

import "fmt"

func main() {
  var functions []func()

  for i := 0; i < 10; i++ {
    j := i
    functions = append(functions, func() {
      fmt.Println(j)
    })
  }

  for _, f := range functions {
    f()
  }
}
```

We could even use `i` as our new variable, so the code could read `i := i` instead of `j := i`. This is called shadowing a variable and can lead to some confusing bugs if abused.

## Summary

+ Functions are first-cass objects in Go. They can be assigned to variables and passed to other functions as parameters.
+ An anonymous function or function literal can be invoked directly after declaration by appending the parameter list to the literal. 
+ A closure is a function value that references variables from outside its body.


## Exercises

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
