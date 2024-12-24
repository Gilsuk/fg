package result_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/gilsuk/fg/result"
	"github.com/stretchr/testify/assert"
)

func TestWrappingFunctionsWithResult(t *testing.T) {
	t.Parallel()

	t.Run("Wrap function that has no argument", func(t *testing.T) {
		t.Parallel()

		returnOK := func() (string, error) {
			return "OK", nil
		}
		var returnResult func() result.Result[string] = result.Wrap(returnOK)
		var resStr result.Result[string] = returnResult()
		called := false

		resStr.Do(func(str string) {
			called = true
			assert.Equal(t, "OK", str)
		}).Catch(func(err error) {
			assert.Fail(t, "Catch should not be called on success")
		})

		assert.True(t, called)
	})

	t.Run("Wrap function that has no argument and returns error", func(t *testing.T) {
		t.Parallel()

		returnErr := func() (string, error) {
			return "", errors.New("dummy error")
		}
		var returnResult func() result.Result[string] = result.Wrap(returnErr)
		var resStr result.Result[string] = returnResult()
		called := false

		resStr.Do(func(s string) {
			assert.Fail(t, "Do should not be called on fail")
		}).Catch(func(err error) {
			called = true
			assert.Error(t, err)
		})

		assert.True(t, called)
	})

	t.Run("test FlatMap for Result[T]", func(t *testing.T) {
		t.Parallel()

		returnErr := func() (string, error) {
			return "78", nil
		}
		strToInt := strconv.Atoi

		var returnResult func() result.Result[string] = result.Wrap(returnErr)

		var resStr result.Result[string] = returnResult()
		var resInt result.Result[int] = result.FlatMap(strToInt)(resStr)

		resInt.Do(func(i int) {
			assert.Equal(t, 78, i)
		}).Catch(func(err error) {
			assert.Fail(t, "Should not be an error")
		})
	})
}
