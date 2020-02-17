/*******************************************************************************
 * Copyright 2020 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package read

import (
	"testing"

	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/delegate"
	dtoBase "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/common/base"
	dtoError "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/common/error"
	dto "github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/dto/v2dot0/core/metadata/addressable/read"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/application/validator"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/infrastructure"
	"github.com/edgexfoundry/edgex-go/internal/pkg/v2/infrastructure/test"

	"github.com/stretchr/testify/assert"
)

// factoryValidRequest returns a valid addressable add request.
func factoryValidRequest(requestID string) *dto.Request {
	return dto.NewRequest(
		dtoBase.NewRequest(requestID),
		test.FactoryRandomString(),
	)
}

// TestValidate tests the Validator for Version request.
func TestValidate(t *testing.T) {
	variations := []*validator.Variation{
		func() *validator.Variation {
			request := factoryValidRequest(test.FactoryRandomString())
			return &validator.Variation{
				Name:             "valid",
				Request:          request,
				ExpectedResponse: request,
				ExpectedStatus:   infrastructure.StatusSuccess,
			}
		}(),
		func() *validator.Variation {
			invalidRequestType := "string is not the request type we're expecting."
			return &validator.Variation{
				Name:    "invalid type",
				Request: invalidRequestType,
				ExpectedResponse: dtoError.NewResponse(
					dtoBase.NewResponse("", invalidRequestType, application.StatusTypeAssertionFailure),
				),
				ExpectedStatus: application.StatusTypeAssertionFailure,
			}
		}(),
		func() *validator.Variation {
			request := factoryValidRequest("")
			return &validator.Variation{
				Name:    "empty requestID",
				Request: request,
				ExpectedResponse: dtoError.NewResponse(
					dtoBase.NewResponse("", request, application.StatusRequestIdEmptyFailure),
				),
				ExpectedStatus: application.StatusRequestIdEmptyFailure,
			}
		}(),
		func() *validator.Variation {
			request := factoryValidRequest(test.FactoryRandomString())
			request.ID = ""

			return &validator.Variation{
				Name:    "missing ID",
				Request: request,
				ExpectedResponse: dtoError.NewResponse(
					dtoBase.NewResponse(request.RequestID, request, application.StatusAddressableMissingID),
				),
				ExpectedStatus: application.StatusAddressableMissingID,
			}
		}(),
	}

	for variationIndex := range variations {
		t.Run(
			variations[variationIndex].Name,
			func(t *testing.T) {
				result, status := Validate(
					variations[variationIndex].Request,
					application.NewBehavior(
						application.Version2,
						application.KindAddressableRead,
						application.ActionRead,
					),
					delegate.TestExecutable,
				)

				assert.Equal(t, variations[variationIndex].ExpectedResponse, result)
				assert.Equal(t, variations[variationIndex].ExpectedStatus, status)
			},
		)
	}
}