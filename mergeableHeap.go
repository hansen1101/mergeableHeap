package mergeableHeap

// A heap Item contains a key used for prioritization and a value which holds the node's data.
type Item interface {
	Key()(int)
	Value()(interface{})
	setKey(int)()
}

// A MergeableHeap is a min heap and provides all operations necessary to implement a priority queue. 
// The elements that are organized in the priority queue must implement the Item interface.
type MergeableHeap interface {
	Insert(Item)(Item)
	Minimum()(Item)
	ExtractMin()(Item)
	Union(MergeableHeap)(MergeableHeap)
}

// A FibonacciHeap is a MergeableHeap that additionally provides DecreaseKey and Delete operations.
type FibonacciHeap interface {
	MergeableHeap
	DecreaseKey(Item,int)()
	Delete(Item)()
}

// Implementatation of the Item interface that has 2 fields of type integer. Thus a node's value is
// simply an integer representing the node's index.
type merableHeapNode struct {
	key int
	value int
}

// Constructor method for creating a merableHeapNode struct.
// @param k integer value for the initial key
// @param v integer value for the node's index
// @return pointer to the allocated merableHeapNode
func NewMergableHeapNode(k,v int)(*merableHeapNode){
	return &merableHeapNode{key:k,value:v}
}

// Implementation of the interface's Key method.
// @return integer value of the node's actual key
func (k merableHeapNode) Key()(int){
	return k.key
}

// Implementation of the interface's Value method.
// @return integer value of the node's actual key
func (k merableHeapNode) Value()(interface{}){
	return k.value
}

// Implementation of the interface's setKey method. This method is package visible and considered not to be accessed from outside.
// The key of any Item can only be changed over the overlying heap implementation.
// @param integer value of the Item's new key
func (k *merableHeapNode) setKey(x int)(){
	k.key = x
}