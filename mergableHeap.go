package mergableHeap

type Item interface {
	Key()(int)
	Value()(int)
	setKey(int)()
}

type merableHeapNode struct {
	key,value int
}

func NewMergableHeapNode(k,v int)(*merableHeapNode){
	return &merableHeapNode{key:k,value:v}
}

func (k merableHeapNode) Key()(int){
	return k.key
}

func (k merableHeapNode) Value()(int){
	return k.value
}

func (k *merableHeapNode) setKey(x int)(){
	k.key = x
}

type MergableHeap interface {
	Insert(Item)(Item)
	ExtractMin()(Item)
	Minimum()(Item)
	DecreaseKey(Item,int)()
	Delete(Item)()
	Union(MergableHeap)(MergableHeap)
}

func MakeHeap()(result *MergableHeap){
	return
}