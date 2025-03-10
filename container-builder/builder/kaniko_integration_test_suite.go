//go:build integration_kaniko_docker

/*
 * Copyright 2023 Red Hat, Inc. and/or its affiliates.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package builder

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/kiegroup/kogito-serverless-operator/container-builder/common"
)

type KanikoDockerTestSuite struct {
	suite.Suite
	LocalRegistry common.DockerLocalRegistry
	RegistryID    string
	Docker        common.Docker
}

func (suite *KanikoDockerTestSuite) SetupSuite() {
	dockerRegistryContainer, registryID, docker := common.SetupDockerSocket()
	if len(registryID) > 0 {
		suite.LocalRegistry = dockerRegistryContainer
		suite.RegistryID = registryID
		suite.Docker = docker
	} else {
		assert.FailNow(suite.T(), "Initialization failed %s", registryID)
	}

	pullErr := suite.Docker.PullImage(EXECUTOR_IMAGE)
	if pullErr != nil {
		logrus.Infof("Pull Kaniko executor Error:%s", pullErr)
	}
	time.Sleep(4 * time.Second) // Needed on CI
}

func (suite *KanikoDockerTestSuite) TearDownSuite() {
	registryID := suite.LocalRegistry.GetRegistryRunningID()
	if len(registryID) > 0 {
		common.DockerTearDown(suite.LocalRegistry)
	} else {
		suite.LocalRegistry.StopRegistry()
	}
	purged, err := suite.Docker.PurgeContainer("", common.REGISTRY_IMG)
	logrus.Infof("Purged containers %t", purged)
	if err != nil {
		logrus.Infof("Purged registry err %t", err)
	}

	purgedBuild, err := suite.Docker.PurgeContainer("", EXECUTOR_IMAGE)
	if err != nil {
		logrus.Infof("Purged container err %t", err)
	}
	logrus.Infof("Purged container build %t", purgedBuild)
}
