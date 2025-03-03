/*
Copyright 2023 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testtools

import (
	"testing"

	"github.com/beego/beego/v2/client/orm"
	"github.com/golang/mock/gomock"

	"github.com/kubeedge/kubeedge/edge/mocks/beego"
	"github.com/kubeedge/kubeedge/edge/pkg/common/dbm"
)
// Comments below is assisted by Gen AI
// // InitOrmerMock initializes and returns mocked instances of beego.Ormer and beego.QuerySeter
// for testing purposes. It sets up the global variables `dbm.DBAccess` and `dbm.DefaultOrmFunc`
// to use the mocked Ormer instance, allowing tests to simulate database interactions without
// requiring a real database connection.
//
// Parameters:
//   - t *testing.T: The testing context provided by the Go testing framework.
//
// Returns:
//   - *beego.MockOrmer: A mocked instance of beego.Ormer.
//   - *beego.MockQuerySeter: A mocked instance of beego.QuerySeter.
//
// Example usage:
//   ormerMock, querySeterMock := InitOrmerMock(t)
//   ormerMock.EXPECT().Insert(gomock.Any()).Return(1, nil)
//   querySeterMock.EXPECT().Filter(gomock.Any(), gomock.Any()).Return(querySeterMock)
func InitOrmerMock(t *testing.T) (*beego.MockOrmer, *beego.MockQuerySeter) {
	//Initialize Global Variables (Mocks)
	// ormerMock is mocked Ormer implementation
	var ormerMock *beego.MockOrmer
	// querySeterMock is mocked QuerySeter implementation
	var querySeterMock *beego.MockQuerySeter

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ormerMock = beego.NewMockOrmer(mockCtrl)
	querySeterMock = beego.NewMockQuerySeter(mockCtrl)
	dbm.DBAccess = ormerMock
	dbm.DefaultOrmFunc = func() orm.Ormer {
		return ormerMock
	}

	return ormerMock, querySeterMock
}
