# PutBucketWebsiteRequestWebsiteConfiguration

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**ErrorDocument** | Pointer to [**ErrorDocument**](ErrorDocument.md) |  | [optional] |
|**IndexDocument** | Pointer to [**IndexDocument**](IndexDocument.md) |  | [optional] |
|**RedirectAllRequestsTo** | Pointer to [**PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo**](PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo.md) |  | [optional] |
|**RoutingRules** | Pointer to [**[]RoutingRule**](RoutingRule.md) |  | [optional] |

## Methods

### NewPutBucketWebsiteRequestWebsiteConfiguration

`func NewPutBucketWebsiteRequestWebsiteConfiguration() *PutBucketWebsiteRequestWebsiteConfiguration`

NewPutBucketWebsiteRequestWebsiteConfiguration instantiates a new PutBucketWebsiteRequestWebsiteConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPutBucketWebsiteRequestWebsiteConfigurationWithDefaults

`func NewPutBucketWebsiteRequestWebsiteConfigurationWithDefaults() *PutBucketWebsiteRequestWebsiteConfiguration`

NewPutBucketWebsiteRequestWebsiteConfigurationWithDefaults instantiates a new PutBucketWebsiteRequestWebsiteConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetErrorDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetErrorDocument() ErrorDocument`

GetErrorDocument returns the ErrorDocument field if non-nil, zero value otherwise.

### GetErrorDocumentOk

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetErrorDocumentOk() (*ErrorDocument, bool)`

GetErrorDocumentOk returns a tuple with the ErrorDocument field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) SetErrorDocument(v ErrorDocument)`

SetErrorDocument sets ErrorDocument field to given value.

### HasErrorDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) HasErrorDocument() bool`

HasErrorDocument returns a boolean if a field has been set.

### GetIndexDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetIndexDocument() IndexDocument`

GetIndexDocument returns the IndexDocument field if non-nil, zero value otherwise.

### GetIndexDocumentOk

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetIndexDocumentOk() (*IndexDocument, bool)`

GetIndexDocumentOk returns a tuple with the IndexDocument field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndexDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) SetIndexDocument(v IndexDocument)`

SetIndexDocument sets IndexDocument field to given value.

### HasIndexDocument

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) HasIndexDocument() bool`

HasIndexDocument returns a boolean if a field has been set.

### GetRedirectAllRequestsTo

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetRedirectAllRequestsTo() PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo`

GetRedirectAllRequestsTo returns the RedirectAllRequestsTo field if non-nil, zero value otherwise.

### GetRedirectAllRequestsToOk

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetRedirectAllRequestsToOk() (*PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo, bool)`

GetRedirectAllRequestsToOk returns a tuple with the RedirectAllRequestsTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectAllRequestsTo

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) SetRedirectAllRequestsTo(v PutBucketWebsiteRequestWebsiteConfigurationRedirectAllRequestsTo)`

SetRedirectAllRequestsTo sets RedirectAllRequestsTo field to given value.

### HasRedirectAllRequestsTo

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) HasRedirectAllRequestsTo() bool`

HasRedirectAllRequestsTo returns a boolean if a field has been set.

### GetRoutingRules

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetRoutingRules() []RoutingRule`

GetRoutingRules returns the RoutingRules field if non-nil, zero value otherwise.

### GetRoutingRulesOk

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) GetRoutingRulesOk() (*[]RoutingRule, bool)`

GetRoutingRulesOk returns a tuple with the RoutingRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingRules

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) SetRoutingRules(v []RoutingRule)`

SetRoutingRules sets RoutingRules field to given value.

### HasRoutingRules

`func (o *PutBucketWebsiteRequestWebsiteConfiguration) HasRoutingRules() bool`

HasRoutingRules returns a boolean if a field has been set.


