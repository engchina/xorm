// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integrations

import (
	"database/sql/driver"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecAndQuery(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UserinfoQuery struct {
		Uid  int
		Name string
	}

	assert.NoError(t, testEngine.Sync2(new(UserinfoQuery)))

	res, err := testEngine.Exec("INSERT INTO "+testEngine.TableName("`userinfo_query`", true)+" (uid, name) VALUES (?, ?)", 1, "user")
	assert.NoError(t, err)
	cnt, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	results, err := testEngine.Query("select * from " + testEngine.TableName("userinfo_query", true))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	id, err := strconv.Atoi(string(results[0]["uid"]))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
	assert.Equal(t, "user", string(results[0]["name"]))
}

func TestExecTime(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UserinfoExecTime struct {
		Uid     int
		Name    string
		Created time.Time
	}

	assert.NoError(t, testEngine.Sync2(new(UserinfoExecTime)))
	now := time.Now()
	res, err := testEngine.Exec("INSERT INTO "+testEngine.TableName("`userinfo_exec_time`", true)+" (uid, name, created) VALUES (?, ?, ?)", 1, "user", now)
	assert.NoError(t, err)
	cnt, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	results, err := testEngine.QueryString("SELECT * FROM " + testEngine.TableName("`userinfo_exec_time`", true))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	assert.EqualValues(t, now.In(testEngine.GetTZLocation()).Format("2006-01-02 15:04:05"), results[0]["created"])

	var uet UserinfoExecTime
	has, err := testEngine.Where("uid=?", 1).Get(&uet)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, now.In(testEngine.GetTZLocation()).Format("2006-01-02 15:04:05"), uet.Created.Format("2006-01-02 15:04:05"))
}

var _ driver.Valuer = &MyData{}

type MyData struct {
	data string
}

func (m *MyData) Value() (driver.Value, error) {
	return m.data, nil
}

func TestExecDriverValuer(t *testing.T) {
	assert.NoError(t, PrepareEngine())

	type UserinfoDriverValuer struct {
		Uid  int
		Name string
		Data string
	}

	assert.NoError(t, testEngine.Sync2(new(UserinfoDriverValuer)))

	res, err := testEngine.Exec("INSERT INTO "+testEngine.TableName("`userinfo_driver_valuer`", true)+" (uid, name,data) VALUES (?, ?, ?)",
		1, "user", MyData{"data"})
	assert.NoError(t, err)
	cnt, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)

	results, err := testEngine.QueryString("select * from " + testEngine.TableName("userinfo_exec", true))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(results))
	id, err := strconv.Atoi(results[0]["uid"])
	assert.NoError(t, err)
	assert.EqualValues(t, 1, id)
	assert.Equal(t, "user", results[0]["name"])
	assert.EqualValues(t, "data", results[0]["data"])
}
