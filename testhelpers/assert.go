package testhelpers

import (
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, name string, value interface{}, expected interface{}) bool {
	t.Helper()
	// if !assertTypeMatch(t,name, value, expected){
	// 	return false
	// }
	//value = upCastBuiltInType(value)
	//expected = upCastBuiltInType(expected)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf(
			"%s[%v] did not euqual expected[%v]",
			name,
			value,
			expected,
		)
		return false
	}
	return true
}
func AssertNotEqual(t *testing.T, name string, value interface{}, expected interface{}) bool {
	t.Helper()
	// if !assertTypeMatch(t,name, value, expected){
	// 	return false
	// }
	//value = upCastBuiltInType(value)
	//expected = upCastBuiltInType(expected)
	if reflect.DeepEqual(value, expected) {
		t.Errorf(
			"%s[%v] did not euqual expected[%v]",
			name,
			value,
			expected,
		)
		return false
	}
	return true
}

// func assertTypeMatch(t *testing.T, name string, value interface{}, expected interface{}) bool {
// 	// i should create a generic switch statement that checks the values of the general types and compares them to what i gave return true if both match
// 	// there is no need for this func yet and I'm unsure on exactly the format of some syntax will finish later

// }

// func upCastBuiltInType(value)
