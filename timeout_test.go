package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeOutFromString_Second_Valid(t *testing.T) {
	timeOut, err := TimeOutFromString("5s")
	assert.Nil(t, err)
	assert.Equal(t, TimeOut{
		Amount:   5,
		TimeUnit: Seconds,
	}, *timeOut)
}

func TestTimeOutFromString_Minutes_Valid(t *testing.T) {
	timeOut, err := TimeOutFromString("5M")
	assert.Nil(t, err)
	assert.Equal(t, TimeOut{
		Amount:   5,
		TimeUnit: Minutes,
	}, *timeOut)
}

func TestTimeOutFromString_Invalid_Format(t *testing.T) {
	timeOut, err := TimeOutFromString("215news23049-34")
	assert.NotNil(t, err)
	assert.Nil(t, timeOut)
}

func TestTimeOutFromString_Invalid_Format_DuplicateUnit(t *testing.T) {
	timeOut, err := TimeOutFromString("215sm")
	assert.NotNil(t, err)
	assert.Nil(t, timeOut)
}

func TestTimeOutFromString_Invalid_Minus(t *testing.T) {
	timeOut, err := TimeOutFromString("-215m")
	assert.NotNil(t, err)
	assert.Nil(t, timeOut)
}

func TestTimeOutFromString_INvalid_NoUnit(t *testing.T) {
	timeOut, err := TimeOutFromString("215a")
	assert.NotNil(t, err)
	assert.Nil(t, timeOut)
}
