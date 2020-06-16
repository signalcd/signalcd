# \DeploymentApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCurrentDeployment**](DeploymentApi.md#GetCurrentDeployment) | **Get** /deployments/current | Get the current Deployment
[**ListDeployments**](DeploymentApi.md#ListDeployments) | **Get** /deployments | List Deployments
[**SetCurrentDeployment**](DeploymentApi.md#SetCurrentDeployment) | **Post** /deployments/current | Set the current Deployment



## GetCurrentDeployment

> Deployment GetCurrentDeployment(ctx, )

Get the current Deployment

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListDeployments

> []Deployment ListDeployments(ctx, )

List Deployments

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetCurrentDeployment

> Deployment SetCurrentDeployment(ctx, optional)

Set the current Deployment

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***SetCurrentDeploymentOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a SetCurrentDeploymentOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **inlineObject** | [**optional.Interface of InlineObject**](InlineObject.md)|  | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)
