package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ebrickdev/ebrick/cache/store"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	store := store.NewMockStore(ctrl)

	// When
	ca := New(store)

	// Then
	assert.IsType(t, new(cache), ca)

	assert.Equal(t, store, ca.GetStore())
}

func TestCacheSet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	value := &struct {
		Hello string
	}{
		Hello: "world",
	}

	mockedStore := store.NewMockStore(ctrl)
	mockedStore.EXPECT().Set(ctx, "my-key", value, store.OptionsMatcher{
		Expiration: 5 * time.Second,
	}).Return(nil)

	cache := New(mockedStore)

	// When
	err := cache.Set(ctx, "my-key", value, store.WithExpiration(5*time.Second))
	assert.Nil(t, err)
}

func TestCacheSetWhenErrorOccurs(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	value := &struct {
		Hello string
	}{
		Hello: "world",
	}

	storeErr := errors.New("an error has occurred while inserting data into store")

	mockedStore := store.NewMockStore(ctrl)
	mockedStore.EXPECT().Set(ctx, "my-key", value, store.OptionsMatcher{
		Expiration: 5 * time.Second,
	}).Return(storeErr)

	cache := New(mockedStore)

	// When
	err := cache.Set(ctx, "my-key", value, store.WithExpiration(5*time.Second))
	assert.Equal(t, storeErr, err)
}

func TestCacheGet(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}

	store := store.NewMockStore(ctrl)
	store.EXPECT().Get(ctx, "my-key").Return(cacheValue, nil)

	cache := New(store)

	// When
	value, err := cache.Get(ctx, "my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
}

func TestCacheGetWhenNotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	returnedErr := errors.New("unable to find item in store")

	store := store.NewMockStore(ctrl)
	store.EXPECT().Get(ctx, "my-key").Return(nil, returnedErr)

	cache := New(store)

	// When
	value, err := cache.Get(ctx, "my-key")

	// Then
	assert.Nil(t, value)
	assert.Equal(t, returnedErr, err)
}

func TestCacheGetWithTTL(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	cacheValue := &struct {
		Hello string
	}{
		Hello: "world",
	}
	expiration := 1 * time.Second

	store := store.NewMockStore(ctrl)
	store.EXPECT().GetWithTTL(ctx, "my-key").
		Return(cacheValue, expiration, nil)

	cache := New(store)

	// When
	value, ttl, err := cache.GetWithTTL(ctx, "my-key")

	// Then
	assert.Nil(t, err)
	assert.Equal(t, cacheValue, value)
	assert.Equal(t, expiration, ttl)
}

func TestCacheGetWithTTLWhenNotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	returnedErr := errors.New("unable to find item in store")
	expiration := 0 * time.Second

	store := store.NewMockStore(ctrl)
	store.EXPECT().GetWithTTL(ctx, "my-key").
		Return(nil, expiration, returnedErr)

	cache := New(store)

	// When
	value, ttl, err := cache.GetWithTTL(ctx, "my-key")

	// Then
	assert.Nil(t, value)
	assert.Equal(t, returnedErr, err)
	assert.Equal(t, expiration, ttl)
}

func TestCacheGetCacheKeyWhenKeyIsString(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	store := store.NewMockStore(ctrl)

	cache := New(store)

	// When
	computedKey := cache.getCacheKey("my-Key")

	// Then
	assert.Equal(t, "my-Key", computedKey)
}

func TestCacheGetCacheKeyWhenKeyIsStruct(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	store := store.NewMockStore(ctrl)

	cache := New(store)

	// When
	key := &struct {
		Hello string
	}{
		Hello: "world",
	}

	computedKey := cache.getCacheKey(key)

	// Then
	assert.Equal(t, "8144fe5310cf0e62ac83fd79c113aad2", computedKey)
}

type StructWithGenerator struct{}

func (_ *StructWithGenerator) GetCacheKey() string {
	return "my-generated-key"
}

func TestCacheGetCacheKeyWhenKeyImplementsGenerator(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	store := store.NewMockStore(ctrl)

	cache := New(store)

	// When
	key := &StructWithGenerator{}

	generatedKey := cache.getCacheKey(key)
	// Then
	assert.Equal(t, "my-generated-key", generatedKey)
}

func TestCacheDelete(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	store := store.NewMockStore(ctrl)
	store.EXPECT().Delete(ctx, "my-key").Return(nil)

	cache := New(store)

	// When
	err := cache.Delete(ctx, "my-key")

	// Then
	assert.Nil(t, err)
}

func TestCacheInvalidate(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	mockedStore := store.NewMockStore(ctrl)
	mockedStore.EXPECT().Invalidate(ctx, store.InvalidateOptionsMatcher{
		Tags: []string{"tag1"},
	}).Return(nil)

	cache := New(mockedStore)

	// When
	err := cache.Invalidate(ctx, store.WithInvalidateTags([]string{"tag1"}))

	// Then
	assert.Nil(t, err)
}

func TestCacheInvalidateWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	expectedErr := errors.New("unexpected error during invalidation")

	mockedStore := store.NewMockStore(ctrl)
	mockedStore.EXPECT().Invalidate(ctx, store.InvalidateOptionsMatcher{
		Tags: []string{"tag1"},
	}).Return(expectedErr)

	cache := New(mockedStore)

	// When
	err := cache.Invalidate(ctx, store.WithInvalidateTags([]string{"tag1"}))

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestCacheClear(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	store := store.NewMockStore(ctrl)
	store.EXPECT().Clear(ctx).Return(nil)

	cache := New(store)

	// When
	err := cache.Clear(ctx)

	// Then
	assert.Nil(t, err)
}

func TestCacheClearWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	expectedErr := errors.New("unexpected error during invalidation")

	store := store.NewMockStore(ctrl)
	store.EXPECT().Clear(ctx).Return(expectedErr)

	cache := New(store)

	// When
	err := cache.Clear(ctx)

	// Then
	assert.Equal(t, expectedErr, err)
}

func TestCacheDeleteWhenError(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	ctx := context.Background()

	expectedErr := errors.New("unable to delete key")

	store := store.NewMockStore(ctrl)
	store.EXPECT().Delete(ctx, "my-key").Return(expectedErr)

	cache := New(store)

	// When
	err := cache.Delete(ctx, "my-key")

	// Then
	assert.Equal(t, expectedErr, err)
}
