package go_client_elastic

import "testing"

func TestGetList(t *testing.T) {
	actual, err := getList(GCERequestBuilder{
		Host:   "10.138.32.97",
		Port:   "8080",
		Index:  "interstellar-gsso-summary",
		Type:   "logs",
		Fields: []string{"message", "AppName"}})
	if &actual == nil {
		t.Errorf("Test failed, nil result(%v).", actual)
	}
	if err != "" {
		t.Errorf("Test failed,Error (%s)", err)
	}
}
func TestStringField(t *testing.T) {
	input := []string{"abc", "def", "ght"}
	expect := "abc,def,ght,"
	actual := stringField(input)
	if actual != expect {
		t.Errorf("Test failed, input(%v), expext(%v), actual(%v).", input, expect, actual)
	}
}
