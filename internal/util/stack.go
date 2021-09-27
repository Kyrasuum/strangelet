package util

type Stack []interface{}

// Create a new stack
func New() *Stack {
	return &Stack{}
}

// Get size of stack
func (this *Stack) Len() int {
	return len(*this)
}

// View the top item on the stack
func (this *Stack) Peek() interface{} {
	if len(*this) == 0 {
		return nil
	}
	return (*this)[0]
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() interface{} {
	if len(*this) == 0 {
		return nil
	}
	elem := this.Peek()
	*this = (*this)[:(len(*this) - 1)]
	return elem
}

// Push a value onto the top of the stack
func (this *Stack) Push(elem interface{}) {
	*this = append(*this, elem)
}
