package form3

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCreateUUID(t *testing.T) {
	v := CreateUUID()
	c := `[0-9a-z]`
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("^%s{8}-%s{4}-%s{4}-%s{4}-%s{12}$", c, c, c, c, c)), v)
}
