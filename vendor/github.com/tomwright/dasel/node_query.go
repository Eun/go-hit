package dasel

import (
	"fmt"
	"reflect"
)

// Query uses the given selector to query the current node and return the result.
func (n *Node) Query(selector string) (*Node, error) {
	n.Selector.Remaining = selector
	rootNode := n

	if err := buildFindChain(rootNode); err != nil {
		return nil, err
	}

	return lastNode(rootNode), nil
}

// lastNode returns the last node in the chain.
// If a node contains multiple next nodes, the first node is taken.
func lastNode(n *Node) *Node {
	node := n
	for {
		if node.Next == nil {
			return node
		}
		node = node.Next
	}
}

func isFinalSelector(selector string) bool {
	return selector == "" || selector == "."
}

func buildFindChain(n *Node) error {
	if isFinalSelector(n.Selector.Remaining) {
		// We've reached the end
		return nil
	}

	var err error
	nextNode := &Node{}

	// Parse the selector.
	nextNode.Selector, err = ParseSelector(n.Selector.Remaining)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %w", err)
	}

	// Link the nodes.
	n.Next = nextNode
	nextNode.Previous = n

	// Populate the value for the new node.
	nextNode.Value, err = findValue(nextNode, false)
	if err != nil {
		return fmt.Errorf("could not find value: %w", err)
	}

	return buildFindChain(nextNode)
}

func findValuePropertyWork(n *Node, createIfNotExists bool, value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
	case reflect.Map:
		for _, key := range value.MapKeys() {
			if fmt.Sprint(key.Interface()) == n.Selector.Property {
				return value.MapIndex(key), nil
			}
		}
		if createIfNotExists {
			return nilValue(), nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}

	case reflect.Struct:
		fieldV := value.FieldByName(n.Selector.Property)
		if fieldV.IsValid() {
			return fieldV, nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}

	case reflect.Ptr:
		return findValuePropertyWork(n, createIfNotExists, derefValue(value))
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// findValueProperty finds the value for the given node using the property selector
// information.
func findValueProperty(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}
	return findValuePropertyWork(n, createIfNotExists, unwrapValue(n.Previous.Value))
}

// findValueIndex finds the value for the given node using the index selector
// information.
func findValueIndex(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	if value.Kind() == reflect.Slice {
		valueLen := value.Len()
		if n.Selector.Index >= 0 && n.Selector.Index < valueLen {
			return value.Index(n.Selector.Index), nil
		}
		if createIfNotExists {
			return nilValue(), nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// findNextAvailableIndex finds the value for the given node using the index selector
// information.
func findNextAvailableIndex(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !createIfNotExists {
		// Next available index isn't supported unless it's creating the item.
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}
	}
	return nilValue(), nil
}

// processFindDynamicItem is used by findValueDynamic.
func processFindDynamicItem(n *Node, object reflect.Value, key string) (bool, error) {
	// Loop through each condition.
	allConditionsMatched := true
	for _, c := range n.Selector.Conditions {
		// If the object doesn't match any checks, return a ValueNotFound.

		var found bool
		var err error
		switch cond := c.(type) {
		case *KeyEqualCondition:
			found, err = cond.Check(reflect.ValueOf(key))
		default:
			found, err = cond.Check(object)
		}

		if err != nil {
			return false, err
		}
		if !found {
			allConditionsMatched = false
			break
		}
	}
	if allConditionsMatched {
		return true, nil
	}
	return false, nil
}

// findValueDynamic finds the value for the given node using the dynamic selector
// information.
func findValueDynamic(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}

	value := unwrapValue(n.Previous.Value)

	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			object := value.Index(i)
			found, err := processFindDynamicItem(n, object, fmt.Sprint(i))
			if err != nil {
				return nilValue(), err
			}
			if found {
				n.Selector.Type = "INDEX"
				n.Selector.Index = i
				return object, nil
			}
		}
		if createIfNotExists {
			n.Selector.Type = "NEXT_AVAILABLE_INDEX"
			return nilValue(), nil
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}

	case reflect.Map:
		for _, key := range value.MapKeys() {
			object := value.MapIndex(key)
			found, err := processFindDynamicItem(n, object, key.String())
			if err != nil {
				return nilValue(), err
			}
			if found {
				n.Selector.Type = "PROPERTY"
				n.Selector.Property = key.String()
				return object, nil
			}
		}
		return nilValue(), &ValueNotFound{Selector: n.Selector.Current, PreviousValue: n.Previous.Value}
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

func findValueLengthWork(n *Node, value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
	case reflect.Slice:
		return reflect.ValueOf(value.Len()), nil

	case reflect.Map:
		return reflect.ValueOf(value.Len()), nil

	case reflect.String:
		return reflect.ValueOf(value.Len()), nil

	case reflect.Struct:
		return reflect.ValueOf(value.NumField()), nil

	case reflect.Ptr:
		return findValueLengthWork(n, derefValue(value))
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// findValueLength returns the length of the current node.
func findValueLength(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}
	return findValueLengthWork(n, unwrapValue(n.Previous.Value))
}

func findValueTypeWork(n *Node, value reflect.Value) (reflect.Value, error) {
	switch value.Kind() {
	case reflect.Slice:
		return reflect.ValueOf("array"), nil

	case reflect.Map:
		return reflect.ValueOf("map"), nil

	case reflect.String:
		return reflect.ValueOf("string"), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf("int"), nil

	case reflect.Float32, reflect.Float64:
		return reflect.ValueOf("float"), nil

	case reflect.Bool:
		return reflect.ValueOf("bool"), nil

	case reflect.Struct:
		return reflect.ValueOf("struct"), nil

	case reflect.Ptr:
		return findValueTypeWork(n, derefValue(value))
	}

	return nilValue(), &UnsupportedTypeForSelector{Selector: n.Selector, Value: value}
}

// findValueType returns the type of the current node.
func findValueType(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if !isValid(n.Previous.Value) {
		return nilValue(), &UnexpectedPreviousNilValue{Selector: n.Previous.Selector.Current}
	}
	return findValueTypeWork(n, unwrapValue(n.Previous.Value))
}

// findValue finds the value for the given node.
// The value is essentially pulled from the previous node, using the (already parsed) selector
// information stored on the current node.
func findValue(n *Node, createIfNotExists bool) (reflect.Value, error) {
	if n.Previous == nil {
		// previous node is required to get it's value.
		return nilValue(), ErrMissingPreviousNode
	}

	if createIfNotExists && !isValid(n.Previous.Value) {
		n.Previous.Value = initialiseEmptyValue(n.Selector, n.Previous.Value)
	}

	switch n.Selector.Type {
	case "PROPERTY":
		return findValueProperty(n, createIfNotExists)
	case "INDEX":
		return findValueIndex(n, createIfNotExists)
	case "NEXT_AVAILABLE_INDEX":
		return findNextAvailableIndex(n, createIfNotExists)
	case "DYNAMIC":
		return findValueDynamic(n, createIfNotExists)
	case "LENGTH":
		return findValueLength(n, createIfNotExists)
	case "TYPE":
		return findValueType(n, createIfNotExists)
	default:
		return nilValue(), &UnsupportedSelector{Selector: n.Selector.Raw}
	}
}
