
# Grammar 

```
digit               =   "0".."9"
startDelimiter = "("
endDelimiter = ")"
delimeter = startDelimiter | endDelimiter
addOperation = "+"
multiplyOperation = "*"
operation = addOperation | multiplyOperation
expression = [startDelimiter],digit|expression,[operation,digit|expression],[endDelimiter]

```
