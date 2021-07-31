package hw04lrucache

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
	list  *list
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func (l list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) Expand() {
	l.length++
}

// добавляем элемент в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	// в конце операции увеличиваем длинну
	defer l.Expand()

	newItem := &ListItem{
		Value: v,
		Next:  &ListItem{},
		Prev:  &ListItem{},
		list:  l,
	}

	// Если первый элемент == nil
	if l.head == nil {
		// присваиваем значения первому элементу
		l.head = newItem
		l.tail = newItem // add
		return newItem
	}
	// меняем текущее значение первому элементу
	currentFront := l.head
	newItem.Next = currentFront
	if currentFront != nil {
		currentFront.Prev = newItem
	}

	l.head = newItem
	return newItem
}

// добавляем элемент в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	defer l.Expand()

	newItem := &ListItem{
		Value: v,
		Next:  &ListItem{},
		Prev:  &ListItem{},
		list:  l,
	}
	// Если последний элемент == nil
	if l.tail == nil {
		// Задаем новое значение первого и последнего элемента
		l.tail = newItem
		l.head = newItem
		return newItem
	}
	// Меняем текущее последнее на получаемое значение.
	currentBack := l.tail
	newItem.Prev = currentBack
	if currentBack != nil {
		currentBack.Next = newItem
	}
	l.tail = newItem
	return newItem
}

// удаляем элемент.
func (l *list) Remove(i *ListItem) {
	// Если элемент один или список пуст, то просто обнуляем список
	if l.length <= 1 {
		l.head = nil
		l.tail = nil
		l.length = 0
		return
	}
	switch i {
	case l.head:
		l.head.Next.Prev = nil
		l.head = l.head.Next
	case l.tail:
		l.tail.Prev.Next = nil
		l.tail = l.tail.Prev
	default:
		// мы знаем куда указывает каждый элемент
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i = nil
	l.length--
}

// премещаем элемент в перед.
func (l *list) MoveToFront(i *ListItem) {
	switch i {
	// не делаем ничего если текущий элемент в начале
	case l.head:
		return
	case l.tail:
		// делаем предпоследний элемент последним
		l.tail = i.Prev
		l.tail.Next = nil
	default:
		// меняем указатели для соседних элементов.
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	// Меняем указатель следующего элемента на текущий и делаем первым текущий
	i.Next = l.head
	l.head, l.head.Prev = i, i
	i.Prev = nil
}

func NewList() List {
	return &list{}
}
