# go-dlx
Exact covering with colors with Go

## example

### nqueen
````
$ go run examples/queen/main.go 4
1:
. Q . . 
. . . Q 
Q . . . 
. . Q . 

2:
. . Q . 
Q . . . 
. . . Q 
. Q . . 

````

### sudoku
````
$ cd examples/sudoku
$ go run main.go < s17.dlx
1:
5 9 7  2 1 8  3 6 4  
1 3 2  4 6 5  8 9 7  
8 6 4  3 7 9  1 2 5  

9 1 5  6 8 4  7 3 2  
3 4 8  7 9 2  6 5 1  
2 7 6  1 5 3  4 8 9  

6 5 9  8 4 7  2 1 3  
4 2 1  9 3 6  5 7 8  
7 8 3  5 2 1  9 4 6  

````
