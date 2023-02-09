/*
 * IONOS Cloud - Managed Stackable Data Platform API
 *
 * Managed Stackable Data Platform by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable Platform.  Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API in addition or as an alternative to the \"Data Center Designer\" (DCD).  ## Getting Started  To get your DataPlatformCluster up and running, the following steps needs to be performed.  ### IONOS Cloud Account  The first step is the creation of a IONOS Cloud account if not already existing.  To register a **new account** visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).  ### Virtual Datacenter (VDC)  The Managed Data Stack needs a virtual datacenter (VDC) hosting the cluster. This could either be a VDC that already exists, especially if you want to connect the managed DataPlatform to other services already running within your VDC. Otherwise, if you want to place the Managed Data Stack in a new VDC or you have not yet created a VDC, you need to do so.  A new VDC can be created via the IONOS Cloud API, the IONOS-CLI or the DCD Web interface. For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/tutorials/data-center-basics)  ### Get a authentication token  To interact with this API a user specific authentication token is needed. This token can be generated using the IONOS-CLI the following way:  ``` ionosctl token generate ```  For more information [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate)  ### Create a new DataPlatformCluster  Before using Managed Stackable Data Platform, a new DataPlatformCluster must be created.  To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.  The provisioning of the cluster might take some time. To check the current provisioning status, you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID that was presented to you in the response of the create cluster call.  ### Add a DataPlatformNodePool  To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool that is provisioned along with the cluster is reserved for the Stackable operators. You may create further node pools with resources tailored to your use-case.  To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post) endpoint.  ### Receive Kubeconfig  Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API. The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.  The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint using the cluster ID of the created DataPlatformCluster.  ### Create Stackable Service  To create the desired application, the Stackable service needs to be provided, using the received kubeconfig and [deploy a Stackable service](https://docs.stackable.tech/home/getting_started.html#_deploying_stackable_services)  ## Authorization  All endpoints are secured, so only an authenticated user can access them. As Authentication mechanism the default IONOS Cloud authentication mechanism is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).  ### Basic-Auth  The basic auth scheme uses the IONOS Cloud user credentials in form of a Basic Authentication Header accordingly to [RFC7617](https://datatracker.ietf.org/doc/html/rfc7617)  ### API-Key as Bearer Token  The Bearer auth token used at the API-Gateway is a user related token created with the IONOS-CLI. (See the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate) for details) For every request to be authenticated, the token is passed as 'Authorization Bearer' header along with the request.  ### Permissions and access roles  Currently, an admin can see and manipulate all resources in a contract. A normal authenticated user can only see and manipulate resources he created.   ## Components  The Managed Stackable Data Platform by IONOS Cloud consists of two components. The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.  ### DataPlatformCluster  A DataPlatformCluster is the virtual instance of the customer services and operations running the managed Services like Stackable operators. A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer. Therefore, it's possible to integrate the cluster with other resources as vLANs e.G. to shape the datacenter in the customer's need and integrate the cluster within the topology the customer wants to build.  In addition to the Kubernetes cluster a small node pool is provided which is exclusively used to run the Stackable operators.  ### DataPlatformNodePool  A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top. All nodes within a node pool are identical in setup. The nodes of a pool are provisioned into virtual data centers at a location of your choice and you can freely specify the properties of all the nodes at once before creation.  Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are readonly in the customer's VDC and can only be modified or deleted via the API.  ### References
 *
 * API version: 0.0.7
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"strings"
	"time"
)

// ToPtr - returns a pointer to the given value.
func ToPtr[T any](v T) *T {
	return &v
}

// ToValue - returns the value of the pointer passed in
func ToValue[T any](ptr *T) T {
	return *ptr
}

// ToValueDefault - returns the value of the pointer passed in, or the default type value if the pointer is nil
func ToValueDefault[T any](ptr *T) T {
	var defaultVal T
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

func SliceToValueDefault[T any](ptrSlice *[]T) []T {
	return append([]T{}, *ptrSlice...)
}

// PtrBool - returns a pointer to given boolean value.
func PtrBool(v bool) *bool { return &v }

// PtrInt - returns a pointer to given integer value.
func PtrInt(v int) *int { return &v }

// PtrInt32 - returns a pointer to given integer value.
func PtrInt32(v int32) *int32 { return &v }

// PtrInt64 - returns a pointer to given integer value.
func PtrInt64(v int64) *int64 { return &v }

// PtrFloat32 - returns a pointer to given float value.
func PtrFloat32(v float32) *float32 { return &v }

// PtrFloat64 - returns a pointer to given float value.
func PtrFloat64(v float64) *float64 { return &v }

// PtrString - returns a pointer to given string value.
func PtrString(v string) *string { return &v }

// PtrTime - returns a pointer to given Time value.
func PtrTime(v time.Time) *time.Time { return &v }

// ToBool - returns the value of the bool pointer passed in
func ToBool(ptr *bool) bool {
	return *ptr
}

// ToBoolDefault - returns the value of the bool pointer passed in, or false if the pointer is nil
func ToBoolDefault(ptr *bool) bool {
	var defaultVal bool
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToBoolSlice - returns a bool slice of the pointer passed in
func ToBoolSlice(ptrSlice *[]bool) []bool {
	valSlice := make([]bool, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToByte - returns the value of the byte pointer passed in
func ToByte(ptr *byte) byte {
	return *ptr
}

// ToByteDefault - returns the value of the byte pointer passed in, or 0 if the pointer is nil
func ToByteDefault(ptr *byte) byte {
	var defaultVal byte
	if ptr == nil {
		return defaultVal
	}

	return *ptr
}

// ToByteSlice - returns a byte slice of the pointer passed in
func ToByteSlice(ptrSlice *[]byte) []byte {
	valSlice := make([]byte, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToString - returns the value of the string pointer passed in
func ToString(ptr *string) string {
	return *ptr
}

// ToStringDefault - returns the value of the string pointer passed in, or "" if the pointer is nil
func ToStringDefault(ptr *string) string {
	var defaultVal string
	if ptr == nil {
		return defaultVal
	}

	return *ptr
}

// ToStringSlice - returns a string slice of the pointer passed in
func ToStringSlice(ptrSlice *[]string) []string {
	valSlice := make([]string, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt - returns the value of the int pointer passed in
func ToInt(ptr *int) int {
	return *ptr
}

// ToIntDefault - returns the value of the int pointer passed in, or 0 if the pointer is nil
func ToIntDefault(ptr *int) int {
	var defaultVal int
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToIntSlice - returns a int slice of the pointer passed in
func ToIntSlice(ptrSlice *[]int) []int {
	valSlice := make([]int, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt8 - returns the value of the int8 pointer passed in
func ToInt8(ptr *int8) int8 {
	return *ptr
}

// ToInt8Default - returns the value of the int8 pointer passed in, or 0 if the pointer is nil
func ToInt8Default(ptr *int8) int8 {
	var defaultVal int8
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt8Slice - returns a int8 slice of the pointer passed in
func ToInt8Slice(ptrSlice *[]int8) []int8 {
	valSlice := make([]int8, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt16 - returns the value of the int16 pointer passed in
func ToInt16(ptr *int16) int16 {
	return *ptr
}

// ToInt16Default - returns the value of the int16 pointer passed in, or 0 if the pointer is nil
func ToInt16Default(ptr *int16) int16 {
	var defaultVal int16
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt16Slice - returns a int16 slice of the pointer passed in
func ToInt16Slice(ptrSlice *[]int16) []int16 {
	valSlice := make([]int16, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt32 - returns the value of the int32 pointer passed in
func ToInt32(ptr *int32) int32 {
	return *ptr
}

// ToInt32Default - returns the value of the int32 pointer passed in, or 0 if the pointer is nil
func ToInt32Default(ptr *int32) int32 {
	var defaultVal int32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt32Slice - returns a int32 slice of the pointer passed in
func ToInt32Slice(ptrSlice *[]int32) []int32 {
	valSlice := make([]int32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt64 - returns the value of the int64 pointer passed in
func ToInt64(ptr *int64) int64 {
	return *ptr
}

// ToInt64Default - returns the value of the int64 pointer passed in, or 0 if the pointer is nil
func ToInt64Default(ptr *int64) int64 {
	var defaultVal int64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt64Slice - returns a int64 slice of the pointer passed in
func ToInt64Slice(ptrSlice *[]int64) []int64 {
	valSlice := make([]int64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint - returns the value of the uint pointer passed in
func ToUint(ptr *uint) uint {
	return *ptr
}

// ToUintDefault - returns the value of the uint pointer passed in, or 0 if the pointer is nil
func ToUintDefault(ptr *uint) uint {
	var defaultVal uint
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUintSlice - returns a uint slice of the pointer passed in
func ToUintSlice(ptrSlice *[]uint) []uint {
	valSlice := make([]uint, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint8 -returns the value of the uint8 pointer passed in
func ToUint8(ptr *uint8) uint8 {
	return *ptr
}

// ToUint8Default - returns the value of the uint8 pointer passed in, or 0 if the pointer is nil
func ToUint8Default(ptr *uint8) uint8 {
	var defaultVal uint8
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint8Slice - returns a uint8 slice of the pointer passed in
func ToUint8Slice(ptrSlice *[]uint8) []uint8 {
	valSlice := make([]uint8, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint16 - returns the value of the uint16 pointer passed in
func ToUint16(ptr *uint16) uint16 {
	return *ptr
}

// ToUint16Default - returns the value of the uint16 pointer passed in, or 0 if the pointer is nil
func ToUint16Default(ptr *uint16) uint16 {
	var defaultVal uint16
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint16Slice - returns a uint16 slice of the pointer passed in
func ToUint16Slice(ptrSlice *[]uint16) []uint16 {
	valSlice := make([]uint16, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint32 - returns the value of the uint32 pointer passed in
func ToUint32(ptr *uint32) uint32 {
	return *ptr
}

// ToUint32Default - returns the value of the uint32 pointer passed in, or 0 if the pointer is nil
func ToUint32Default(ptr *uint32) uint32 {
	var defaultVal uint32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint32Slice - returns a uint32 slice of the pointer passed in
func ToUint32Slice(ptrSlice *[]uint32) []uint32 {
	valSlice := make([]uint32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint64 - returns the value of the uint64 pointer passed in
func ToUint64(ptr *uint64) uint64 {
	return *ptr
}

// ToUint64Default - returns the value of the uint64 pointer passed in, or 0 if the pointer is nil
func ToUint64Default(ptr *uint64) uint64 {
	var defaultVal uint64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint64Slice - returns a uint63 slice of the pointer passed in
func ToUint64Slice(ptrSlice *[]uint64) []uint64 {
	valSlice := make([]uint64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToFloat32 - returns the value of the float32 pointer passed in
func ToFloat32(ptr *float32) float32 {
	return *ptr
}

// ToFloat32Default - returns the value of the float32 pointer passed in, or 0 if the pointer is nil
func ToFloat32Default(ptr *float32) float32 {
	var defaultVal float32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToFloat32Slice - returns a float32 slice of the pointer passed in
func ToFloat32Slice(ptrSlice *[]float32) []float32 {
	valSlice := make([]float32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToFloat64 - returns the value of the float64 pointer passed in
func ToFloat64(ptr *float64) float64 {
	return *ptr
}

// ToFloat64Default - returns the value of the float64 pointer passed in, or 0 if the pointer is nil
func ToFloat64Default(ptr *float64) float64 {
	var defaultVal float64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToFloat64Slice - returns a float64 slice of the pointer passed in
func ToFloat64Slice(ptrSlice *[]float64) []float64 {
	valSlice := make([]float64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToTime - returns the value of the Time pointer passed in
func ToTime(ptr *time.Time) time.Time {
	return *ptr
}

// ToTimeDefault - returns the value of the Time pointer passed in, or 0001-01-01 00:00:00 +0000 UTC if the pointer is nil
func ToTimeDefault(ptr *time.Time) time.Time {
	var defaultVal time.Time
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToTimeSlice - returns a Time slice of the pointer passed in
func ToTimeSlice(ptrSlice *[]time.Time) []time.Time {
	valSlice := make([]time.Time, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

type NullableBool struct {
	value *bool
	isSet bool
}

func (v NullableBool) Get() *bool {
	return v.value
}

func (v *NullableBool) Set(val *bool) {
	v.value = val
	v.isSet = true
}

func (v NullableBool) IsSet() bool {
	return v.isSet
}

func (v *NullableBool) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBool(val *bool) *NullableBool {
	return &NullableBool{value: val, isSet: true}
}

func (v NullableBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBool) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt struct {
	value *int
	isSet bool
}

func (v NullableInt) Get() *int {
	return v.value
}

func (v *NullableInt) Set(val *int) {
	v.value = val
	v.isSet = true
}

func (v NullableInt) IsSet() bool {
	return v.isSet
}

func (v *NullableInt) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt(val *int) *NullableInt {
	return &NullableInt{value: val, isSet: true}
}

func (v NullableInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt32 struct {
	value *int32
	isSet bool
}

func (v NullableInt32) Get() *int32 {
	return v.value
}

func (v *NullableInt32) Set(val *int32) {
	v.value = val
	v.isSet = true
}

func (v NullableInt32) IsSet() bool {
	return v.isSet
}

func (v *NullableInt32) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt32(val *int32) *NullableInt32 {
	return &NullableInt32{value: val, isSet: true}
}

func (v NullableInt32) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt32) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt64 struct {
	value *int64
	isSet bool
}

func (v NullableInt64) Get() *int64 {
	return v.value
}

func (v *NullableInt64) Set(val *int64) {
	v.value = val
	v.isSet = true
}

func (v NullableInt64) IsSet() bool {
	return v.isSet
}

func (v *NullableInt64) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt64(val *int64) *NullableInt64 {
	return &NullableInt64{value: val, isSet: true}
}

func (v NullableInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt64) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableFloat32 struct {
	value *float32
	isSet bool
}

func (v NullableFloat32) Get() *float32 {
	return v.value
}

func (v *NullableFloat32) Set(val *float32) {
	v.value = val
	v.isSet = true
}

func (v NullableFloat32) IsSet() bool {
	return v.isSet
}

func (v *NullableFloat32) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFloat32(val *float32) *NullableFloat32 {
	return &NullableFloat32{value: val, isSet: true}
}

func (v NullableFloat32) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFloat32) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableFloat64 struct {
	value *float64
	isSet bool
}

func (v NullableFloat64) Get() *float64 {
	return v.value
}

func (v *NullableFloat64) Set(val *float64) {
	v.value = val
	v.isSet = true
}

func (v NullableFloat64) IsSet() bool {
	return v.isSet
}

func (v *NullableFloat64) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFloat64(val *float64) *NullableFloat64 {
	return &NullableFloat64{value: val, isSet: true}
}

func (v NullableFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFloat64) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableString struct {
	value *string
	isSet bool
}

func (v NullableString) Get() *string {
	return v.value
}

func (v *NullableString) Set(val *string) {
	v.value = val
	v.isSet = true
}

func (v NullableString) IsSet() bool {
	return v.isSet
}

func (v *NullableString) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableString(val *string) *NullableString {
	return &NullableString{value: val, isSet: true}
}

func (v NullableString) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableString) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableTime struct {
	value *time.Time
	isSet bool
}

func (v NullableTime) Get() *time.Time {
	return v.value
}

func (v *NullableTime) Set(val *time.Time) {
	v.value = val
	v.isSet = true
}

func (v NullableTime) IsSet() bool {
	return v.isSet
}

func (v *NullableTime) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTime(val *time.Time) *NullableTime {
	return &NullableTime{value: val, isSet: true}
}

func (v NullableTime) MarshalJSON() ([]byte, error) {
	return v.value.MarshalJSON()
}

func (v *NullableTime) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type IonosTime struct {
	time.Time
}

func (t *IonosTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if strlen(str) == 0 {
		t = nil
		return nil
	}
	if str[0] == '"' {
		str = str[1:]
	}
	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}
	if !strings.Contains(str, "Z") {
		/* forcefully adding timezone suffix to be able to parse the
		 * string using RFC3339 */
		str += "Z"
	}
	tt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = IonosTime{tt}
	return nil
}
