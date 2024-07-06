// Common tools and helper functions
package common

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/go-playground/validator/v10"
	"github.com/s2dio-tech/mindgra-backend/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// A helper function to generate random string
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// My own Error type that will help return my customized Error info
//
//	{"database": {"hello":"no such table", error: "not_exists"}}
type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

// To handle the error returned by c.Bind in gin framework
// https://github.com/go-playground/validator/blob/v9/_examples/translations/main.go
func NewValidatorError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		// can translate each error one at a time.
		//fmt.Println("gg",v.NameNamespace)
		if v.Param() != "" {
			res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag(), v.Param())
		} else {
			res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag())
		}

	}
	return res
}

// Warp the error info in a object
func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func ConvertToDistinctArray(edges []domain.WordsLink) []string {
	distinctMap := make(map[string]bool)
	var output []string

	for _, edge := range edges {
		distinctMap[edge.SourceId] = true
		distinctMap[edge.TargetId] = true
	}

	for key := range distinctMap {
		output = append(output, key)
	}

	return output
}

func ConvertMapToStruct(m map[string]interface{}, s interface{}) error {
	// Convert the map to JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Convert the JSON to a struct
	json.Unmarshal(jsonData, &s)
	return nil
}
