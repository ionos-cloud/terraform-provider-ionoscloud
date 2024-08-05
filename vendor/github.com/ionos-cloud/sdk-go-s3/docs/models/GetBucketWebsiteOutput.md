# GetBucketWebsiteOutput

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**RedirectAllRequestsTo** | Pointer to [**RedirectAllRequestsTo**](RedirectAllRequestsTo.md) |  | [optional] |
|**IndexDocument** | Pointer to [**IndexDocument**](IndexDocument.md) |  | [optional] |
|**ErrorDocument** | Pointer to [**ErrorDocument**](ErrorDocument.md) |  | [optional] |
|**RoutingRules** | Pointer to [**[]RoutingRule**](RoutingRule.md) |  | [optional] |

## Methods

### NewGetBucketWebsiteOutput

`func NewGetBucketWebsiteOutput() *GetBucketWebsiteOutput`

NewGetBucketWebsiteOutput instantiates a new GetBucketWebsiteOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetBucketWebsiteOutputWithDefaults

`func NewGetBucketWebsiteOutputWithDefaults() *GetBucketWebsiteOutput`

NewGetBucketWebsiteOutputWithDefaults instantiates a new GetBucketWebsiteOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRedirectAllRequestsTo

`func (o *GetBucketWebsiteOutput) GetRedirectAllRequestsTo() RedirectAllRequestsTo`

GetRedirectAllRequestsTo returns the RedirectAllRequestsTo field if non-nil, zero value otherwise.

### GetRedirectAllRequestsToOk

`func (o *GetBucketWebsiteOutput) GetRedirectAllRequestsToOk() (*RedirectAllRequestsTo, bool)`

GetRedirectAllRequestsToOk returns a tuple with the RedirectAllRequestsTo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectAllRequestsTo

`func (o *GetBucketWebsiteOutput) SetRedirectAllRequestsTo(v RedirectAllRequestsTo)`

SetRedirectAllRequestsTo sets RedirectAllRequestsTo field to given value.

### HasRedirectAllRequestsTo

`func (o *GetBucketWebsiteOutput) HasRedirectAllRequestsTo() bool`

HasRedirectAllRequestsTo returns a boolean if a field has been set.

### GetIndexDocument

`func (o *GetBucketWebsiteOutput) GetIndexDocument() IndexDocument`

GetIndexDocument returns the IndexDocument field if non-nil, zero value otherwise.

### GetIndexDocumentOk

`func (o *GetBucketWebsiteOutput) GetIndexDocumentOk() (*IndexDocument, bool)`

GetIndexDocumentOk returns a tuple with the IndexDocument field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndexDocument

`func (o *GetBucketWebsiteOutput) SetIndexDocument(v IndexDocument)`

SetIndexDocument sets IndexDocument field to given value.

### HasIndexDocument

`func (o *GetBucketWebsiteOutput) HasIndexDocument() bool`

HasIndexDocument returns a boolean if a field has been set.

### GetErrorDocument

`func (o *GetBucketWebsiteOutput) GetErrorDocument() ErrorDocument`

GetErrorDocument returns the ErrorDocument field if non-nil, zero value otherwise.

### GetErrorDocumentOk

`func (o *GetBucketWebsiteOutput) GetErrorDocumentOk() (*ErrorDocument, bool)`

GetErrorDocumentOk returns a tuple with the ErrorDocument field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrorDocument

`func (o *GetBucketWebsiteOutput) SetErrorDocument(v ErrorDocument)`

SetErrorDocument sets ErrorDocument field to given value.

### HasErrorDocument

`func (o *GetBucketWebsiteOutput) HasErrorDocument() bool`

HasErrorDocument returns a boolean if a field has been set.

### GetRoutingRules

`func (o *GetBucketWebsiteOutput) GetRoutingRules() []RoutingRule`

GetRoutingRules returns the RoutingRules field if non-nil, zero value otherwise.

### GetRoutingRulesOk

`func (o *GetBucketWebsiteOutput) GetRoutingRulesOk() (*[]RoutingRule, bool)`

GetRoutingRulesOk returns a tuple with the RoutingRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoutingRules

`func (o *GetBucketWebsiteOutput) SetRoutingRules(v []RoutingRule)`

SetRoutingRules sets RoutingRules field to given value.

### HasRoutingRules

`func (o *GetBucketWebsiteOutput) HasRoutingRules() bool`

HasRoutingRules returns a boolean if a field has been set.


