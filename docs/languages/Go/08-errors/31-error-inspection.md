---
id: error_inspection
title: Error inspection
sidebar_label: Error inspection
sidebar_position: 31
draft: true
---

In the previous lecture, we learned about wrapping errors using the `%w` verb in an `fmt.Errorf()` format `string`. To recap, this line

```go
return fmt.Errorf("propagate(%d): %w", i, err)
```

wraps a new error with additional contextual information around an existing error, and passes both back to the caller. It is good practice to add context information when passing an error returned by a function call on to the own caller.

Let me stress the importance of this. Without adding context, an error would simply bubble up the call chain until some caller addresses the error. This caller might need to know why and how the originating error arrived at this point. When every logical layer adds information about the circumstances of the originating error, the code that finally handles this error can better decide what to do.

So how can the error handler inspect the error and all of its wrapped errors, in order to treat the error appropriately?

Go provides three functions for easy error inspection: `Unwrap()`, `Is()`, and `As()`. All three are included in the errors package.

Let's have a look at an example to see them in action. Our sample code shall write some structured data to a file. 

Function `WriteDoc()` calls `os.Open()` to open a file. If that fails, it takes the error returned by `os.Open()` and wraps its own error message around it. The new message contains document title and ID as new context information.

```go
type Doc struct {
    ID    int
    Title string
    Text  string
}

func WriteDoc(path string, doc Doc) error {
    f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0660)
    if err != nil {
        return fmt.Errorf("Cannot write %s (id %d): %w", doc.Title, doc.ID, err)
    }
    // write data...
    return nil
}

func main() {
    doc := Doc{
        ID: 20,
        Title: "Error Inspection",
        Text: "In the previous lecture, we learned about wrapping errors...",
    }

    err := WriteDoc("/path/to/no_file", doc)
    if err != nil {
        fmt.Println("Top-level error:", err) 
    }
}
```

Now the calling function can do a few cool things with the returned error.

## Unwrapping the inner error

For logging or other purposes, we can extract the inner error(s) by calling the library function `Unwrap()`:

```go
    unwrapped := errors.Unwrap(err)
    fmt.Println("Unwrapped error:", unwrapped) 
```

`func Unwrap(err error) error` takes an error value and returns the error inside that error, or `nil` if there is no error inside the given error to unwrap.

Note that `Unwrap()` is not a method of `err` but rather a function of the errors package. You thus need to import errors in order to unwrap an error. The same applies to the two other package functions discussed below.

Let's have a look at the output. The top-level error message contains the error string as defined by the `Errorf()` call above. It includes the wrapped error message:

