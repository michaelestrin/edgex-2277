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

package infrastructure

// Status is the common type definition for the service's transport-agnostic status code.
//
// - infrastructure layer: 1 - 9999
// - domain layer: 10000 - 19999
// - application layer: 20000 - 29999
// - user interface layer: 30000 - 39999
type Status int

const (
	StatusSuccess             Status = 0
	StatusPersistenceNotFound Status = 1
)
