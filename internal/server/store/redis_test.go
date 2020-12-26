package store

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestSetList(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	type testCase struct {
		name      string
		key       string
		ttl       int
		valuesExp []interface{}
		values    []interface{}
		isError   bool
	}

	tCases := []testCase{
		testCase{
			name: "Success",
			key:  "user:1",
			ttl:  10,
			valuesExp: []interface{}{
				`{"Dtype":"string","Data":"Ivan"}`,
				`{"Dtype":"float64","Data":"195.50000"}`,
				`{"Dtype":"int","Data":"21"}`,
				`{"Dtype":"map","Data":"{\"name\":\"Ivan\"}"}`,
			},
			values: []interface{}{
				"Ivan",
				195.5,
				21,
				map[string]interface{}{
					"name": "Ivan",
				},
			},
			isError: false,
		},
		{
			name:    "No fields",
			key:     "user:2",
			ttl:     0,
			isError: true,
		},
		{
			name: "No key",
			ttl:  10,
			valuesExp: []interface{}{
				`{"Dtype":"string","Data":"Ivan"}`,
			},
			values: []interface{}{
				"Ivan",
			},
			isError: true,
		},
		{
			name: "Unvalid ttl",
			key:  "user:1",
			ttl:  -10,
			valuesExp: []interface{}{
				`{"Dtype":"string","Data":"Ivan"}`,
				`{"Dtype":"float64","Data":195.5}`,
				`{"age":"int64","Data":21}`,
			},
			values: []interface{}{
				"Ivan",
				195.5,
				21,
			},
			isError: true,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectRPush(tc.key, tc.valuesExp...).SetVal(3)
			err := client.SetList(context.Background(), tc.key, tc.values, 0)

			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetHash(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	type testCase struct {
		name      string
		key       string
		ttl       int
		valuesExp map[string]interface{}
		values    map[string]interface{}
		isError   bool
	}

	tCases := []testCase{
		testCase{
			name: "Success",
			key:  "user:1",
			ttl:  10,
			valuesExp: map[string]interface{}{
				"name": "Ivan",
				"age":  21,
			},
			values: map[string]interface{}{
				"name": "Ivan",
				"age":  21,
			},
			isError: false,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectHMSet(tc.key, tc.valuesExp).SetVal(true)
			err := client.SetHash(context.Background(), tc.key, tc.values, 0)

			if tc.isError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetString(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "ivan"
	valuesExp := "Lapshin"

	mock.ExpectSet(key, valuesExp, 0).SetVal("OK")
	_, err := client.SetString(context.Background(), key, valuesExp, 0)
	assert.NoError(t, err)
}

func TestGetKeys(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	pattern := "*"
	mock.ExpectKeys(pattern).SetVal([]string{"one", "two", "three"})
	_, err := client.GetKeys(context.Background(), pattern)
	assert.NoError(t, err)
}

func TestDel(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "Ivan"
	mock.ExpectDel(key).SetVal(1)
	_, err := client.Delete(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetString(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	value := "Ivan"
	mock.ExpectGet(key).SetVal(value)
	_, err := client.GetString(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetHash(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	value := map[string]string{
		"name":     "Ivan",
		"lastname": "lapshin",
	}
	mock.ExpectHGetAll(key).SetVal(value)
	_, err := client.GetHash(context.Background(), key)
	assert.NoError(t, err)
}

func TestGetList(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	values := []string{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}

	mock.ExpectLRange(key, 0, -1).SetVal(values)
	_, err := client.GetList(context.Background(), key)
	assert.NoError(t, err)
}

func TestHGet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	field := "name"
	mock.ExpectHGet(key, field).SetVal("user:1")

	_, err := client.HGet(context.Background(), key, field)
	assert.NoError(t, err)
}

func TestHSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "user:1"
	value := map[string]interface{}{
		"name":     "Ivan",
		"lastname": "Lapshin",
		"age":      "21",
	}

	mock.ExpectHSet(key, value).SetVal(3)

	_, err := client.HSet(context.Background(), key, value)
	assert.NoError(t, err)
}

func TestLRange(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "list:1"
	values := []string{
		`{"Dtype":"string","Data":"1"}`, `{"Dtype":"string","Data":"2"}`, `{"Dtype":"string","Data":"qwerty"}`,
	}

	mock.ExpectLRange(key, 0, 1).SetVal(values)
	_, err := client.LRange(context.Background(), key, 0, 1)
	assert.NoError(t, err)
}

func TestLSet(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := Redis{
		client: db,
	}

	key := "list:1"
	value := `{"Dtype":"string","Data":"1"}`
	index := 1

	mock.ExpectLSet(key, int64(index), value)

	_, err := client.LSet(context.Background(), key, int64(index), value)

	assert.NoError(t, err)
}
