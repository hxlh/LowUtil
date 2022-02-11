package bplustree

type KeyType int

const kM int = 3
const kMINKEY_NUM int = ((kM + 2 - 1) / 2) - 1
const kINVALID_KEY int = ^int(^uint(0) >> 1)

type BPlusNode struct {
	isLeaf   bool
	parent   *BPlusNode
	next     *BPlusNode
	key_nums int
	keys     [kM]KeyType
	children [kM + 1]*BPlusNode
}

func (n *BPlusNode) keyIndex(key KeyType) int {
	//二分查找
	l := 0
	r := n.key_nums
	for l < r {
		mid := (l + r) / 2
		if n.keys[mid] > key {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return r
}

func (n *BPlusNode) childIndex(child *BPlusNode) int {
	for i := 0; i < n.key_nums+1; i++ {
		if n.children[i] == child {
			return i
		}
	}
	return -1
}

func (n *BPlusNode) insertKey(key KeyType) {
	index := n.keyIndex(key)
	for i := n.key_nums - 1; i >= index; i-- {
		n.keys[i+1] = n.keys[i]
	}
	n.keys[index] = key
	n.key_nums++
}

func (n *BPlusNode) deleteKey(key KeyType) {
	index := n.keyIndex(key)
	if n.keys[index-1] == key {
		for i := index - 1; i < n.key_nums; i++ {
			n.keys[i] = n.keys[i+1]
			n.keys[i+1] = KeyType(kINVALID_KEY)
		}
		n.key_nums--
	}
}

func (n *BPlusNode) findSibling() (lSibling int, rSibling int) {
	if n.parent == nil {
		return -1, -1
	}
	parent := n.parent

	for i := 0; i < parent.key_nums+1; i++ {
		if parent.children[i] == n {
			return i - 1, i + 1
		}
	}

	return -1, -1
}

func (n *BPlusNode) deleteNode(node *BPlusNode) bool {
	i := 0
	for ; i < n.key_nums+1; i++ {
		if n.children[i] == node {
			break
		}
	}
	if i == n.key_nums+1 {
		return false
	}
	for i = i + 1; i < n.key_nums+1; i++ {
		n.children[i-1] = n.children[i]
		n.children[i] = nil
	}
	return true
}

func (n *BPlusNode) insertNodeWithIndex(index int, node *BPlusNode) {
	for i := n.key_nums; i >= index; i-- {
		n.children[i+1] = n.children[i]
		n.children[i] = nil
	}
	n.children[index] = node
	node.parent = n
}

func (n *BPlusNode) appendNode(node *BPlusNode) {
	if node == nil {
		return
	}
	for i := 0; i < len(n.children); i++ {
		if n.children[i] == nil {
			n.children[i] = node
			node.parent = n
			break
		}
	}
}

func (n *BPlusNode) insertNode(oldNode *BPlusNode, newNode *BPlusNode) {

	index := 0
	for ; index < n.key_nums+1; index++ {
		if oldNode == n.children[index] {
			break
		}
	}

	//找不到oldNode
	if index == n.key_nums+1 {
		//根节点无children情况
		n.children[0] = oldNode
		n.children[1] = newNode
		oldNode.parent = n
		newNode.parent = n
	} else if index == n.key_nums {
		//oldNode 为最后一个节点
		n.children[n.key_nums] = newNode
		newNode.parent = n
	} else {
		for i := n.key_nums - 1; i >= index+1; i-- {
			n.children[i+1] = n.children[i]
			n.children[i] = nil
		}
		n.children[index+1] = newNode
		newNode.parent = n
	}

	//维系关系
	if oldNode.isLeaf {
		newNode.next = oldNode.next
		oldNode.next = newNode
	}
}

type BPlusTree struct {
	root *BPlusNode
}

func (b *BPlusTree) Insert(key KeyType) {
	//根节点为空
	if b.root == nil {
		b.root = &BPlusNode{}
		b.root.isLeaf = true

		b.root.insertKey(key)
		return
	}

	cursor := b.root
	//find leaf
	for !cursor.isLeaf {
		index := cursor.keyIndex(key)
		cursor = cursor.children[index]
	}

	cursor.insertKey(key)

	if cursor.key_nums < kM {
		return
	}
	//需要分裂
	for cursor.parent != nil {
		parent := cursor.parent
		mid := cursor.key_nums / 2
		parent.insertKey(cursor.keys[mid])

		newNode := &BPlusNode{}
		b.initNode(newNode)
		newNode.isLeaf = cursor.isLeaf
		newNode.parent = parent
		if cursor.isLeaf {
			//拷贝key
			for i := mid; i < cursor.key_nums; i++ {
				newNode.insertKey(cursor.keys[i])
			}
			//拷贝children
			for i := mid + 1; i < cursor.key_nums+1; i++ {
				newNode.appendNode(cursor.children[i])
				cursor.children[i] = nil
			}
		} else {
			//拷贝key
			for i := mid + 1; i < cursor.key_nums; i++ {
				newNode.insertKey(cursor.keys[i])
			}
			//拷贝children
			for i := mid + 1; i < cursor.key_nums+1; i++ {
				newNode.appendNode(cursor.children[i])
				cursor.children[i] = nil
			}
		}
		cursor.key_nums = mid

		parent.insertNode(cursor, newNode)
		cursor = cursor.parent
		if cursor.key_nums < kM {
			return
		}
	}
	//cursor.key_nums>=kM && cursor==root
	root := &BPlusNode{}
	newNode := &BPlusNode{}
	b.initNode(root)
	b.initNode(newNode)
	newNode.isLeaf = cursor.isLeaf
	newNode.parent = root
	cursor.parent = root
	if cursor.isLeaf {
		newNode.next = cursor.next
		cursor.next = newNode
	}
	mid := cursor.key_nums / 2
	root.insertKey(cursor.keys[mid])
	if cursor.isLeaf {
		//拷贝key
		for i := mid; i < cursor.key_nums; i++ {
			newNode.insertKey(cursor.keys[i])
		}
		//拷贝children
		for i := mid + 1; i < cursor.key_nums+1; i++ {
			newNode.appendNode(cursor.children[i])
			cursor.children[i] = nil
		}
	} else {
		//拷贝key
		for i := mid + 1; i < cursor.key_nums; i++ {
			newNode.insertKey(cursor.keys[i])
		}
		//拷贝children
		for i := mid + 1; i < cursor.key_nums+1; i++ {
			newNode.appendNode(cursor.children[i])
			cursor.children[i] = nil
		}
	}
	cursor.key_nums = mid

	root.appendNode(cursor)
	root.appendNode(newNode)
	b.root = root
	b.root.parent = nil
}

func (b *BPlusTree) Delete(key KeyType) bool {
	cursor := b.root
	var internalNode *BPlusNode
	for !cursor.isLeaf {
		index := cursor.keyIndex(key)
		if index != 0 && cursor.keys[index-1] == key {
			internalNode = cursor
		}
		cursor = cursor.children[index]
	}
	index := cursor.keyIndex(key)
	//找不到key
	if index <= 0 || cursor.keys[index-1] != key {
		return false
	}

	cursor.deleteKey(key)

	parent := cursor.parent
	if cursor == b.root && cursor.key_nums <= 0 {
		b.root = nil
		return true
	} else if internalNode == nil {
		//索引没有这个key(keynums必然大于kmin)
		if cursor.key_nums >= kMINKEY_NUM {
			return true
		}
		//需要找叶子兄弟借一个key
		leftSibIndex, rightSibIndex := cursor.findSibling()
		if leftSibIndex >= 0 {
			leftSib := parent.children[leftSibIndex]
			if leftSib.key_nums > kMINKEY_NUM {
				parent.deleteKey(cursor.keys[0])
				parent.insertKey(leftSib.keys[leftSib.key_nums-1])
				cursor.insertKey(leftSib.keys[leftSib.key_nums-1])
				leftSib.deleteKey(leftSib.keys[leftSib.key_nums-1])
				return true
			}
			if rightSibIndex > 0 && rightSibIndex <= parent.key_nums {
				rightSib := parent.children[rightSibIndex]
				if rightSib.key_nums > kMINKEY_NUM {
					parent.deleteKey(rightSib.keys[0])
					cursor.insertKey(rightSib.keys[0])
					rightSib.deleteKey(rightSib.keys[0])
					parent.insertKey(rightSib.keys[0])
					return true
				}
			}
		}
		//借不到key,跟兄弟节点合并
		//如果左边有兄弟则跟左边合并
		if leftSibIndex >= 0 {
			leftSib := parent.children[leftSibIndex]
			parent.deleteNode(cursor)
			parent.deleteKey(cursor.keys[0])
			for i := 0; i < cursor.key_nums; i++ {
				leftSib.insertKey(cursor.keys[i])
			}
		} else {
			rightSib := parent.children[rightSibIndex]
			parent.deleteNode(cursor)
			parent.deleteKey(rightSib.keys[0])
			for i := 0; i < cursor.key_nums; i++ {
				rightSib.insertKey(cursor.keys[i])
			}

		}
		//合并后判断parent是否满足b+tree规则
		b.internalAdjust(parent)
		return true
	} else if internalNode == parent {
		//向左右兄弟借一个key
		index := parent.childIndex(cursor)
		leftSib := parent.children[index-1]
		if leftSib.key_nums > kMINKEY_NUM {
			cursor.insertKey(leftSib.keys[leftSib.key_nums-1])
			parent.deleteKey(parent.keys[index-1])
			parent.insertKey(leftSib.keys[leftSib.key_nums-1])
			leftSib.deleteKey(leftSib.keys[leftSib.key_nums-1])
			return true
		} else if index+1 <= parent.key_nums {
			rightSib := parent.children[index+1]
			if rightSib.key_nums > kMINKEY_NUM {
				parent.deleteKey(key)
				cursor.insertKey(rightSib.keys[0])
				rightSib.deleteKey(rightSib.keys[0])
				parent.insertKey(rightSib.keys[0])
				return true
			}
		}

		//借不到
		//合并key
		for i := 0; i < cursor.key_nums; i++ {
			leftSib.insertKey(cursor.keys[i])
		}

		parent.deleteNode(cursor)
		parent.deleteKey(parent.keys[index-1])

		b.internalAdjust(parent)
		return true
	} else if internalNode != parent {
		//向右兄弟借一个key
		index := parent.childIndex(cursor)
		rightSib := parent.children[index+1]
		if rightSib.key_nums > kMINKEY_NUM {
			cursor.insertKey(rightSib.keys[0])
			parent.deleteKey(rightSib.keys[0])
			rightSib.deleteKey(rightSib.keys[0])
			parent.insertKey(rightSib.keys[0])
			internalNode.deleteKey(key)
			internalNode.insertKey(cursor.keys[0])
			return true
		}
		//跟右兄弟合并
		parent.deleteNode(cursor)
		internalNode.deleteKey(key)
		parent.deleteKey(rightSib.keys[0])
		//合并key
		for i := 0; i < cursor.key_nums; i++ {
			rightSib.insertKey(cursor.keys[0])
		}
		internalNode.insertKey(rightSib.keys[0])
		b.internalAdjust(parent)
		return true
	}
	return false
}

func (b *BPlusTree) internalAdjust(n *BPlusNode) {
	cursor := n
	if cursor.key_nums >= kMINKEY_NUM {
		return
	}
	for cursor.parent != nil {
		parent := cursor.parent
		leftSibIndex, rightSibIndex := cursor.findSibling()
		//先向兄弟借一个key
		if leftSibIndex >= 0 {
			leftSib := parent.children[leftSibIndex]
			if leftSib.key_nums > kMINKEY_NUM {
				index := parent.childIndex(cursor)
				cursor.insertKey(parent.keys[index-1])
				cursor.insertNodeWithIndex(0, leftSib.children[leftSib.key_nums])
				parent.deleteKey(parent.keys[index-1])
				parent.insertKey(leftSib.keys[leftSib.key_nums-1])
				leftSib.deleteNode(leftSib.children[leftSib.key_nums])
				leftSib.deleteKey(leftSib.keys[leftSib.key_nums-1])
				return
			}
			if rightSibIndex > 0 && rightSibIndex <= parent.key_nums {
				rightSib := parent.children[rightSibIndex]
				if rightSib.key_nums > kMINKEY_NUM {
					index := parent.childIndex(cursor)
					cursor.insertKey(parent.keys[index])
					cursor.insertNodeWithIndex(cursor.key_nums+1, rightSib.children[0])
					parent.deleteKey(parent.keys[0])
					parent.insertKey(rightSib.keys[0])
					rightSib.deleteNode(rightSib.children[0])
					rightSib.deleteKey(rightSib.keys[0])
					return
				}
			}
		}
		//向左右兄弟合并
		if leftSibIndex >= 0 {
			leftSib := parent.children[leftSibIndex]
			index := leftSibIndex + 1
			leftSib.insertKey(parent.keys[index-1])

			for i := 0; i < cursor.key_nums+1; i++ {
				leftSib.insertNodeWithIndex(leftSib.key_nums, cursor.children[i])
			}
			for i := 0; i < cursor.key_nums; i++ {
				leftSib.insertKey(cursor.keys[i])
			}
			parent.deleteNode(cursor)
			parent.deleteKey(parent.keys[index-1])
		} else {
			rightSib := parent.children[rightSibIndex]
			index := rightSibIndex - 1
			rightSib.insertKey(parent.keys[index])

			for i := cursor.key_nums; i >= 0; i-- {
				rightSib.insertNodeWithIndex(0, cursor.children[i])
			}
			for i := 0; i < cursor.key_nums; i++ {
				rightSib.insertKey(cursor.keys[i])
			}
			parent.deleteNode(cursor)
			parent.deleteKey(parent.keys[index])
		}
		if parent.key_nums >= kMINKEY_NUM {
			return
		}
		cursor = parent
	}
	if cursor.key_nums >= 1 {
		return
	}
	b.root = cursor.children[0]
	b.root.parent = nil
}

func (b *BPlusTree) initNode(node *BPlusNode) {
	for i := 0; i < len(node.children); i++ {
		node.children[i] = nil
	}
	for i := 0; i < len(node.keys); i++ {
		node.keys[i] = KeyType(kINVALID_KEY)
	}
	node.key_nums = 0
	node.parent = nil
	node.isLeaf = false
}

func (b *BPlusTree) Display(ctl string) {
	println("\n----------" + ctl + "-------------")
	q := make([]*BPlusNode, 0)
	q = append(q, b.root)
	q = append(q, &BPlusNode{key_nums: kM})
	var node *BPlusNode
	for len(q) > 0 {
		node = q[0]
		q = q[1:]
		if node == nil {
			break
		}
		if node.key_nums == kM {
			println("")
			continue
		}

		for i := 0; i < node.key_nums+1; i++ {
			q = append(q, node.children[i])
		}
		if q[0].key_nums >= kM {
			q = append(q, &BPlusNode{key_nums: kM})
		}

		for i := 0; i < node.key_nums; i++ {
			print(node.keys[i])
			print("-")
		}
		print("	")
	}
}
