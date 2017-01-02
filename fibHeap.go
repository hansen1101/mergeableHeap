package mergableHeap

import(
	"container/list"
	"fmt"
	"bytes"
	"math"
)

type FibHeapNode struct {
	Item
	Parent *FibHeapNode
	Child *list.List
	mark bool
}

func NewFibHeapNode(k,v int)(*FibHeapNode){
	return &FibHeapNode{Item:NewMergableHeapNode(k,v)}
}

func (n FibHeapNode) Rank()(int){
	return n.Child.Len()
}

type FibonacciHeap struct {
	MergableHeap
	rootlist *list.List
	minimum *FibHeapNode
	n int // number of nodes in heap
}

// MergableHeap interface methods

// Works
func (h *FibonacciHeap) Insert(z Item)(Item){
	var x *FibHeapNode
	var ok bool
	x,ok = z.(*FibHeapNode)
	if !ok {
		x = NewFibHeapNode(z.Key(),z.Value())
	}
	x.Child = list.New()
	x.Parent = nil
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

// works
func (h *FibonacciHeap) ExtractMin()(Item){
	
	var z *FibHeapNode
	z = h.minimum
	
	//fmt.Printf("Extract Min Start Root list list size:%v, n:%v \t",h.rootlist.Len(),h.n)

	if z != nil {
		
		for e := h.rootlist.Front(); e != nil; e = e.Next() {
			node := e.Value.(*FibHeapNode)
			if (*z).Item == (*node).Item {
				h.rootlist.Remove(e)
				h.n--
				break
			}
		}

		// iterate over children of b-tree root
		for x := z.Child.Front(); x != nil; x = x.Next() {
			//add x to the root list
			child := x.Value.(*FibHeapNode)
			child.Parent = nil
			h.rootlist.PushBack(child)
		}

		//fmt.Printf("Extracted Minimum Key: %d\tAfter Min Extraction list size:%v, n:%v\n",z.Key(),h.rootlist.Len(),h.n)
		
		if h.rootlist.Len()==0 {
			h.minimum = nil
		} else {
			min := h.rootlist.Front().Value.(*FibHeapNode)
			h.minimum = min
			h.Consolidate()
		}
	}
	return z.Item
}

func (h *FibonacciHeap) Minimum()(Item){
	return h.minimum
}

// Works
func (h *FibonacciHeap) DecreaseKey(z Item,key int)(){
	x,ok := z.(*FibHeapNode)
	//@todo: try to find a suitable FibHeapNode for the Item if no FibHeapNode pointer is provided
	if !ok || x.Key() < key {
		//fmt.Println("Decrease Key failed.")
		return
	}
	x.setKey(key)
	y := x.Parent
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

func (h *FibonacciHeap) Delete(x Item)(){
	return
}

// Works
func (h1 *FibonacciHeap) Union(z MergableHeap)(MergableHeap){
	h2,ok := z.(*FibonacciHeap)
	if !ok {
		//@todo check if there is a way to merge non fibheap with h1
		return nil
	}
	h := &FibonacciHeap{}
	h.rootlist = h1.rootlist
	h.minimum = h1.minimum
	h.rootlist.PushBackList(h2.rootlist)
	switch {
	case h.minimum == nil:
		h.minimum = h2.minimum
	case h2.minimum != nil:
		if h2.minimum.Key() < h.minimum.Key() {
			h.minimum = h2.minimum
		}
	}
	h.n = h1.n + h2.n
	return h
}

// FibonacciHeap specific methods

func FibMakeHeap()(result *FibonacciHeap){
	result = &FibonacciHeap{}
	result.rootlist = list.New()
	return
}

func (h *FibonacciHeap) Len()int{
	return h.n
}

func (h FibonacciHeap) Degree()(int){
	return h.rootlist.Len()
}

func (h *FibonacciHeap) Cut(x,y *FibHeapNode)(){
	// remove x from the cild list of y, decrementing y.degree
	for element := y.Child.Front(); element != nil; element = element.Next(){
		if element.Value.(*FibHeapNode) == x {
			y.Child.Remove(element)
			break
		}
	}
	x.Parent = nil
	x.mark = false
	h.rootlist.PushBack(x)
	return
}

func (h *FibonacciHeap) CascadeCut(y *FibHeapNode)(){
	z := y.Parent
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
func (h *FibonacciHeap) D()(maxDegree int){
	d := 1.0 / math.Log2((1.0+math.Sqrt(5.0))/2.0) * math.Log2(float64(h.n))
	return int(d)
	/*
	for root := h.rootlist.Front(); root != nil; root = root.Next() {
		if count := root.Value.(*FibHeapNode).Child.Len(); count > maxDegree {
			maxDegree = count
		}
	}
	return
	*/
}

// works
func (h *FibonacciHeap) FibHeapLink(y,x *list.Element)(){
	//remove y from root list of h
	temp := h.rootlist.Remove(y).(*FibHeapNode)
	temp.Parent = x.Value.(*FibHeapNode)
	temp.mark = false
	x.Value.(*FibHeapNode).Child.PushFront(temp)
	return
}

// works
func (h *FibonacciHeap) Consolidate()(){
	deg := h.D()+10
	var a []*list.Element

	if deg > 0 {
		a = make([]*list.Element,deg,deg)
	}
	for i,_ := range a {
		a[i] = nil
	}

	var w, cand *list.Element
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

func (h FibonacciHeap) String() string {
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
	if x.Child.Len()==0 {
		return fmt.Sprintf("(%d,%d,%d,%v)",x.Key(),x.Value(),x.Child.Len(),x.mark)
	} else {
		buf := new(bytes.Buffer)
		buf.WriteString(fmt.Sprintf("(%d,%d,%d,%v - [",x.Key(),x.Value(),x.Child.Len(),x.mark))
		for e := x.Child.Front(); e!=nil; e = e.Next() {
			next := e.Value.(*FibHeapNode)
			buf.WriteString(next.String())
		}
		buf.WriteString("])")
		return buf.String()
	}
}