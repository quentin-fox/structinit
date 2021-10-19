# structinit

`structinit` is a static analysis tool for Go that helps to identify uninitialized values in specific structs. This can be used to ensure that all the fields in a specific struct literal have been set - this can be useful when copying fields from one struct to another, or to make sure all the dependencies for a struct have been initialized.

## Installation

```sh
go get -u github.com/quentin-fox/structinit/cmd/structinit
```

## Usage

```sh
structinit [package]
```

## Tagging Structs

`structinit` will only validate the fields of structs which have been declared using the `var` keyword and tagged with the `structinit:exhaustive` comment on the line immediately before the declaration. Note that adding a tag before a declaration with the `:=` operator will not work.

```go
type Cat struct {
  Name string
  Color string
  Floofiness int
  Friendly bool
}

//structinit:exhaustive
var cat = Cat{ // fails with "Exhaustive struct literal Cat not initialized with field Friendly"
  Name: "Bad Kitty",
  Color: "Calico",
  Floofiness: 6,
}
```

## Omitting Fields from Validation

It is also possible to omit any number of fields from validation, so the struct will be considered exhaustively initialized even without the omitted fields. Using the previous example:

```go
type Cat struct {
  Name string
  Color string
  Floofiness int
  Friendly bool
}

//structinit:exhaustive,omit=Friendly
var cat = Cat{ // no errors reported
  Name: "Bad Kitty",
  Color: "Calico",
  Floofiness: 6,
}
```

Any number of fields can be omitted by passing them as a comma-separated list to `omit`:

```go
//structinit:exhaustive,omit=Floofiness,Friendly
```

`structinit` will not report an error if an omitted field is initialized. However, it will report an error if an omitted field is not one of the fields of the struct being validated.
