package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInterest(t *testing.T) {
	app, cleanUp := setup(context.Background())
	defer cleanUp()

	asserts := assert.New(t)

	asserts.NotNil(app)
}
