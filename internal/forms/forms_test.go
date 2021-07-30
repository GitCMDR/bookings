package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/random-link", nil)
	form := New(r.PostForm)

	isValid := form.Valid()

	if !isValid {
		t.Error("Got invalid when it should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a","b","c") // check if form has fields a, b, c

	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("Shows does not have required field when it does")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil) // create a request

	postedData := url.Values{}
	postedData.Add("name", "luis") // add field name with value luis to form
	r.PostForm = postedData // append data to PostForm

	form := New(r.PostForm) // use the request to create a form

	if form.Has("dog", r) {
		t.Error("Form does not have field 'dog' but function Has found it")
	}

	if !form.Has("name", r) {
		t.Error("Form has field 'name' but function did not find it")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil) // create a request

	postedData := url.Values{}
	postedData.Add("name", "luis") // add field name with value luis to form
	r.PostForm = postedData // append data to PostForm
	form := New(r.PostForm) // use the request to create a form

	if !form.MinLength("name", 4, r) {
		t.Error("Value in field is long enough but validation says it's not")
	}

	if form.MinLength("name", 5, r) {
		t.Error("Value in field is not long enough but validation says it is")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil) // create a request

	postedData := url.Values{}
	postedData.Add("first_email", "luis@luis.com")
	postedData.Add("second_email", "luis@lu")// add field name with value luis to form
	r.PostForm = postedData // append data to PostForm
	form := New(r.PostForm) // use the request to create a form

	form.IsEmail("first_email")
	if !form.Valid() {
		t.Error("Email is valid but validation failed")
	}

	isError := form.Errors.Get("first_email")
	if isError != "" {
		t.Error("Shouldn't have an error but got one")
	}

	form.IsEmail("second_email")
	if form.Valid() {
		t.Error("Email is not valid but validation succeeded")
	}

	isError = form.Errors.Get("second_email")
	if isError == "" {
		t.Error("Should have an error but didn't get one")
	}

}