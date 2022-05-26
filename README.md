# Monkey

Simple implementation Monkey Programming Language while 
reading [Writing An Interpreter In Go](https://interpreterbook.com/).

## Build

```bash
make
```

## Examples

### Variables

```
let age = 1;
let name = "Monkey";
let result = 10 * (20 / 2);
```

### Arrays & Hashes

```
let myArray = [1, 2, 3, 4, 5];
let thorsten = {"name": "Thorsten", "age": 28};
let myArray2 = ["Thorsten", "Ball", 28, fn(x) { x * x }];

myArray[0] // => 1 
thorsten["name"] // => "Thorsten"
```

### Functions

```
let add = fn(a, b) { return a + b; };
let add = fn(a, b) { a + b; };

add(1, 2);

let fibonacci = fn(x) { 
    if (x == 0) {
        0
    } else {
        if (x == 1) {
            1
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        } 
    }
};

let twice = fn(f, x) { 
    return f(f(x));
};

let addTwo = fn(x) { 
    return x + 2;
};
twice(addTwo, 2); // => 6
```

### Builtin functions
```
let myArray = ["one", "two", "three"];

len(myArray)            // 3
first(myArray)          // one
rest(myArray)           // [two, three]
last(myArray)           // three
push(myArray, "four")   // [one, two, three, four]

len("Hey Bob, how ya doin?") // 21

puts("Hello World!")
```

## TODOs: For the future (maybe)

- [ ] Implement UTF-8 support
- [ ] Implement float numbers
- [ ] Implement reading code from files
- [ ] Implement postfix expressions
