package mergeableHeap

import(
	"container/list"
	"fmt"
	"bytes"
	"math"
)

type FibHeapNode struct {
	Item				// is an Item and therefore all methods of the underlying Item are available
	parent *FibHeapNode // pointer to the parent of a node or nil if the node is a root and therefore in the rootlist
	child *list.List    // list storing all children of the node
	mark bool			// marker field shows if one of the node's children was deleted since the last time the node was added to a child list of another node
}

// Constructor method to generate a new FibHeapNode struct. The method assigns the value und key given as paramenters
// and returns a pointer to the newly allocated FibHeapNode struct.
func NewFibHeapNode(k,v int)(*FibHeapNode){
	return &FibHeapNode{Item:NewMergableHeapNode(k,v)}
}

// The rank or degree of a FibHeapNode is the number of its children.
// @return int the number of the child nodes stored in the child list.
func (n FibHeapNode) Rank()(int){
	return n.child.Len()
}

type FibonacciHeapImpl struct {
	FibonacciHeap           // a FibonacciHeapImpl is a FibonacciHeap
	rootlist *list.List     // pointer to the rootlist containting all root nodes
	minimum *FibHeapNode    // pointer to the node containing the minimum key
	n int 					// number of nodes in heap
}

// MergeableHeap interface methods

// Inserts an Item into the rootlist of the heap. If the Item is not a FibHeapNode, a new FibHeapNode is allocated.
// @param z Item that should be inserted into the heap's rootlist
// @return pointer to the Item that is inserted.
func (h *FibonacciHeapImpl) Insert(z Item)(Item){
	var x *FibHeapNode
	var ok bool
	x,ok = z.(*FibHeapNode)
	if !ok {
		x = NewFibHeapNode(z.Key(),z.Value().(int))
	}
	x.child = list.New()
	x.parent = nil
	x.mark = false
	if h.minimum == nil {
		h.minimum = x
	}
	h.rootlist.PushBack(x)
	h.n++
	if x.Key() < h.minimum.Key() {
		h.minimum = x
	}
	return x
}

// Returns a pointer to the Item with the minimum key in the heap.
func (h *FibonacciHeapImpl) Minimum()(Item){
	return h.minimum
}
// Finds and deletes the Item with the minimum key from the heap. The method returns a pointer to the deleted Item.
// After the minimum is deleted it is also ensured that the fibonacci heap structure is restored properly.
// @return pointer to the Item with the minimum key
func (h *FibonacciHeapImpl) ExtractMin()(Item){
	var z *FibHeapNode
	z = h.minimum
	if z != nil {
		// loop over rootlist and delete list element holding pointer reference to the current minimum z
		for e := h.rootlist.Front(); e != nil; e = e.Next() {
			node := e.Value.(*FibHeapNode)
			if (*z).Item == (*node).Item {
				h.rootlist.Remove(e)
				h.n--
				break
			}
		}
		// loop over z's children and add them to the rootlist
		for x := z.child.Front(); x != nil; x = x.Next() {
			child := x.Value.(*FibHeapNode)
			child.parent = nil
			h.rootlist.PushBack(child)
		}
		if h.rootlist.Len()==0 {
			h.minimum = nil
		} else {
			// set minimum pointer to front of the rootlist and consolidate the rootlist elements
			h.minimum = h.rootlist.Front().Value.(*FibHeapNode)
			h.Consolidate()
		}
	}
	return z.Item
}

// Merges the heaps h and z into one heap.
// @param z pointer to a MergeableHeap that should be merged with h.
// @return pointer to a MergeableHeap that consists of both h and z.
func (h *FibonacciHeapImpl) Union(z MergeableHeap)(MergeableHeap){
	h2,ok := z.(*FibonacciHeapImpl)
	if !ok {
		//@todo check if there is a way to merge non fibheap with h
		return nil
	}
	y := &FibonacciHeapImpl{}
	y.rootlist = h.rootlist
	y.minimum = h.minimum
	y.rootlist.PushBackList(h2.rootlist)
	switch {
	case y.minimum == nil:
		y.minimum = h2.minimum
	case h2.minimum != nil:
		if h2.minimum.Key() < y.minimum.Key() {
			y.minimum = h2.minimum
		}
	}
	y.n = h.n + h2.n
	return y
}

// FibonacciHeap interface methods

// Works
func (h *FibonacciHeapImpl) DecreaseKey(z Item,key int)(){
	x,ok := z.(*FibHeapNode)
	//@todo: try to find a suitable FibHeapNode for the Item if no FibHeapNode pointer is provided
	if !ok || x.Key() < key {
		//fmt.Println("Decrease Key failed.")
		return
	}
	x.setKey(key)
	y := x.parent
	if y != nil {
		if x.Key() < y.Key() {
			h.Cut(x,y)
			h.CascadeCut(y)
		}
	}
	if h.minimum == nil {
		h.minimum = x
	}
	if x.Key() < h.minimum.Key() {
		h.minimum = x
	}
	return
}


func (h *FibonacciHeapImpl) Delete(x Item)(){
	h.DecreaseKey(x,math.MinInt32)
	h.ExtractMin()
	return
}


// FibonacciHeapImpl specific methods

func FibMakeHeap()(result *FibonacciHeapImpl){
	result = &FibonacciHeapImpl{}
	result.rootlist = list.New()
	return
}

func (h *FibonacciHeapImpl) Len()int{
	return h.n
}

func (h FibonacciHeapImpl) Degree()(int){
	return h.rootlist.Len()
}

func (h *FibonacciHeapImpl) Cut(x,y *FibHeapNode)(){
	// remove x from the cild list of y, decrementing y.degree
	for element := y.child.Front(); element != nil; element = element.Next(){
		if element.Value.(*FibHeapNode) == x {
			y.child.Remove(element)
			break
		}
	}
	x.parent = nil
	x.mark = false
	h.rootlist.PushBack(x)
	return
}

func (h *FibonacciHeapImpl) CascadeCut(y *FibHeapNode)(){
	z := y.parent
	if z != nil {
		if !y.mark {
			y.mark = true
		} else {
			h.Cut(y,z)
			h.CascadeCut(z)
		}
	}
	return
}

// Works
func (h *FibonacciHeapImpl) D()(maxDegree int){
	d := 1.0 / math.Log2((1.0+math.Sqrt(5.0))/2.0) * math.Log2(float64(h.n)) + 1.0
	return int(d)
}

// works
func (h *FibonacciHeapImpl) FibHeapLink(y,x *list.Element)(){
	//remove y from root list of h
	temp := h.rootlist.Remove(y).(*FibHeapNode)
	temp.parent = x.Value.(*FibHeapNode)
	temp.mark = false
	x.Value.(*FibHeapNode).child.PushFront(temp)
	return
}

// works
func (h *FibonacciHeapImpl) Consolidate()(){
	var a []*list.Element
	var w, cand *list.Element
	deg := h.D()
	if deg > 0 {
		a = make([]*list.Element,deg,deg)
	}
	for i,_ := range a {
		a[i] = nil
	}
	w = h.rootlist.Front()
	for w!=nil {
		cand = w.Next()
		x:=w
		d := x.Value.(*FibHeapNode).Rank()
		for a[d] != nil {
			y := a[d]
			if x.Value.(*FibHeapNode).Key() > y.Value.(*FibHeapNode).Key() {
				temp:=x
				x = y
				y = temp
			}
			h.FibHeapLink(y,x)
			a[d] = nil
			d++
			if d >= len(a){
				fmt.Printf("%v - %v\n",d,len(a))
			}
		}
		a[d] = x
		w = cand
	}
	h.minimum = nil
	for _,v := range a {
		if v != nil {
			if h.minimum == nil {
				h.minimum = v.Value.(*FibHeapNode)
				newrootlist := list.New()
				h.rootlist = newrootlist
				h.rootlist.PushFront(h.minimum)
			} else {
				if v.Value.(*FibHeapNode).Key() < h.minimum.Key() {
					h.minimum = v.Value.(*FibHeapNode)
					h.rootlist.PushFront(h.minimum)
				} else {
					h.rootlist.PushBack(v.Value.(*FibHeapNode))
				}
			}
		}
	}
	return
}

func (h FibonacciHeapImpl) String() string {
	buf := new(bytes.Buffer)
	var x *list.Element
	for e := h.rootlist.Front(); e!=nil; e = e.Next() {
		if e.Value.(*FibHeapNode) == h.minimum {
			x = e
			break
		}
	}

	buf.WriteString(fmt.Sprintf("Total Elements: %d\nRoot Elements: %d\nMinimum Key: %d\n",h.n,h.rootlist.Len(),h.minimum.Key()))
	buf.WriteString(fmt.Sprintf("Min:\t%s\n",x.Value.(*FibHeapNode)))

	for next := x.Next(); next != nil; next = next.Next(){
		buf.WriteString(fmt.Sprintf("\t %s\n",next.Value.(*FibHeapNode)))
	}

	for next := x.Prev(); next != nil; next = next.Prev(){
		buf.WriteString(fmt.Sprintf("\t %s\n",next.Value.(*FibHeapNode)))
	}

	return buf.String()
}

func (x *FibHeapNode) String() string {
	if x.child.Len()==0 {
		return fmt.Sprintf("(%d,%d,%d,%v)",x.Key(),x.Value(),x.child.Len(),x.mark)
	} else {
		buf := new(bytes.Buffer)
		buf.WriteString(fmt.Sprintf("(%d,%d,%d,%v - [",x.Key(),x.Value(),x.child.Len(),x.mark))
		for e := x.child.Front(); e!=nil; e = e.Next() {
			next := e.Value.(*FibHeapNode)
			buf.WriteString(next.String())
		}
		buf.WriteString("])")
		return buf.String()
	}
}