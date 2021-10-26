package main

import "bufio"
import "fmt"
import "os"

import "example.com/stack"

// Global operator and operand stacks
var operandStack stack.Stack
var operatorStack stack.Stack

// Returns x ^ y. This is a brute force integer power routine using successive
// multiplication. (There are more efficient ways to do this.)
func intPower(x int, y int) (pow int) {
	pow = 1
	for i := 0 ; i < y ; i++ {
		pow *= x
	}
	return
}

// Returns true if the character is a digit.
func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// Returns the precedence of the operator.
func precedence(op byte) (prec int) {
	switch op {
	case '+', '-': prec = 0
	case '*', '/': prec = 1
	case '^': prec = 2
	default: panic("unknown operator")
	}
	return
}

// Returns true if op is right associative. Only exponentiation is right
// associative.
func isRightAssociative(op byte) bool {
    return op == '^'
}

// Apply the top operator on the operator stack to the top two operands on the
// operand stand and push the result onto the operand stack.
func apply() {
	// Pop the operator off the operator stack
	op, err := operatorStack.Pop()
	if err != nil {
		panic("operator stack underflow")
	}
	// Pop the right operand off the operand stack
	right, err := operandStack.Pop()
	if err != nil {
		panic("operand stack underflow")
	}
	// Pop the left operand off the operand stack
	left, err := operandStack.Pop()
	if err != nil {
		panic("operand stack underflow")
	}
	// Apply the operator to the left and right operands and push the result
	// onto the operand stack
	switch op.(byte) {
	case '+': operandStack.Push(left.(int) + right.(int))
	case '-': operandStack.Push(left.(int) - right.(int))
	case '*': operandStack.Push(left.(int) * right.(int))
	case '/': operandStack.Push(left.(int) / right.(int))
	case '^': operandStack.Push(intPower(left.(int), right.(int)))
	default: panic("unknown operator")
	}
	return
}

// Evaluate an expression and print the result.
func evaluate(expr string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("illegal expression:", r)
		}
	}()
	
	// Process the expression character by character left to right
	operandExpected := true
	i := 0
	for i < len(expr) {
		switch expr[i] {
		// Digit: Extract the operand and push it on the operand stack
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if !operandExpected {
				panic("operator expected but operand found")
			}
			v := 0
			for i < len(expr) && isDigit(expr[i]) {
				v = 10*v+int(expr[i]-'0')
				i += 1
			}
			operandStack.Push(v)
			operandExpected = false
		// Operator: Apply pending operators of greater or equal precedence
		// then push the operator on the operator stack
		case '+', '-', '*', '/', '^':
			if operandExpected {
				panic("operand expected but operator found")
			}
			for !operatorStack.IsEmpty() {
				op, _ := operatorStack.Top()
				if precedence(op.(byte)) > precedence(expr[i]) ||
				   (precedence(op.(byte)) == precedence(expr[i]) &&
				    !isRightAssociative(op.(byte))) {
					apply()
				} else {
					break
				}
			}
			operatorStack.Push(expr[i])
			i += 1
			operandExpected = true
		case '(':
			fmt.Printf("%q is an open parenthesis\n", expr[i])
			i += 1
		case ')':
			fmt.Printf("%q is a close parenthesis\n", expr[i])
			i += 1
		case ' ':
			i += 1
		default:
			panic(fmt.Sprintf("%q is an illegal character", expr[i]))
		}
	}
	// Apply any remaining operators
	for !operatorStack.IsEmpty() {
		apply()
	}
	// The result is the one operator remaining on the stack.
	result, _ := operandStack.Pop()
	if !operandStack.IsEmpty() {
		panic("too many operands")
	}
	fmt.Printf("%v\n", result)
}

// Main routine to read expressions from standard input, calculate their values,
// and print the result. (Use an end of file, control-Z, to exit.)
func main() {
	
	// Make a scanner to read lines from standard input
	scanner := bufio.NewScanner(os.Stdin)
	
	// Process each of the lines from standard input
	for scanner.Scan() {
	
		// Initialize the operator and operand stacks
		operandStack = stack.New()
		operatorStack = stack.New()
		
		// Get the current line of text.
		line := scanner.Text()
		// fmt.Println(line)
		
		// Evaluate the expression and print the result
		evaluate(line)
	}
}
