package rbd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupSnapshots(t *testing.T) {
	// tests are done as subtests to avoid creating pools, images, etc
	// over and over again.
	conn := radosConnect(t)
	require.NotNil(t, conn)
	defer conn.Shutdown()

	poolname := GetUUID()
	err := conn.MakePool(poolname)
	require.NoError(t, err)
	defer conn.DeletePool(poolname)

	ioctx, err := conn.OpenIOContext(poolname)
	require.NoError(t, err)
	defer ioctx.Destroy()

	// create a group, some images, and add images to the group
	gname := "snapme"
	err = GroupCreate(ioctx, gname)
	assert.NoError(t, err)

	options := NewRbdImageOptions()
	assert.NoError(t,
		options.SetUint64(ImageOptionOrder, uint64(testImageOrder)))

	name1 := GetUUID()
	err = CreateImage(ioctx, name1, testImageSize, options)
	require.NoError(t, err)

	name2 := GetUUID()
	err = CreateImage(ioctx, name2, testImageSize, options)
	require.NoError(t, err)

	err = GroupImageAdd(ioctx, gname, ioctx, name1)
	assert.NoError(t, err)
	err = GroupImageAdd(ioctx, gname, ioctx, name2)
	assert.NoError(t, err)

	t.Run("groupSnapCreateRemove", func(t *testing.T) {
		err := GroupSnapCreate(ioctx, gname, "snap1")
		assert.NoError(t, err)
		err = GroupSnapRemove(ioctx, gname, "snap1")
		assert.NoError(t, err)
	})
	t.Run("groupSnapRename", func(t *testing.T) {
		err := GroupSnapCreate(ioctx, gname, "snap2a")
		assert.NoError(t, err)
		err = GroupSnapRename(ioctx, gname, "fred", "wilma")
		assert.Error(t, err)
		err = GroupSnapRename(ioctx, gname, "snap2a", "snap2b")
		assert.NoError(t, err)
		err = GroupSnapRemove(ioctx, gname, "snap2a")
		assert.Error(t, err, "remove of old name: expect error")
		err = GroupSnapRemove(ioctx, gname, "snap2b")
		assert.NoError(t, err, "remove of current name: expect success")
	})
	t.Run("invalidIOContext", func(t *testing.T) {
		assert.Panics(t, func() {
			GroupSnapCreate(nil, gname, "foo")
		})
		assert.Panics(t, func() {
			GroupSnapRemove(nil, gname, "foo")
		})
		assert.Panics(t, func() {
			GroupSnapRename(nil, gname, "foo", "bar")
		})
	})
}