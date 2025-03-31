package twerge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArbitraryShadow(t *testing.T) {
	assert.Equal(t, true, isArbitraryShadow("[inset_0_1px_0,inset_0_-1px_0]"))
	assert.Equal(t, true, isArbitraryShadow("[0_35px_60px_-15px_rgba(0,0,0,0.3)]"))
	assert.Equal(t, true, isArbitraryShadow("[inset_0_1px_0,inset_0_-1px_0]"))
	assert.Equal(t, true, isArbitraryShadow("[0_0_#00f]"))
	assert.Equal(t, true, isArbitraryShadow("[.5rem_0_rgba(5,5,5,5)]"))
	assert.Equal(t, true, isArbitraryShadow("[-.5rem_0_#123456]"))
	assert.Equal(t, true, isArbitraryShadow("[0.5rem_-0_#123456]"))
	assert.Equal(t, true, isArbitraryShadow("[0.5rem_-0.005vh_#123456]"))
	assert.Equal(t, true, isArbitraryShadow("[0.5rem_-0.005vh]"))

	assert.Equal(t, false, isArbitraryShadow("[rgba(5,5,5,5)]"))
	assert.Equal(t, false, isArbitraryShadow("[#00f]"))
	assert.Equal(t, false, isArbitraryShadow("[something-else]"))
}
