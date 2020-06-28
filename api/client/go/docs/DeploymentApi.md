# \DeploymentApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCurrentDeployment**](DeploymentApi.md#GetCurrentDeployment) | **Get** /deployments/current | Get the current Deployment
[**ListDeployments**](DeploymentApi.md#ListDeployments) | **Get** /deployments | List Deployments
[**SetCurrentDeployment**](DeploymentApi.md#SetCurrentDeployment) | **Post** /deployments/current | Set the current Deployment
[**UpdateDeploymentStatus**](DeploymentApi.md#UpdateDeploymentStatus) | **Patch** /deployments/{id}/status | Update parts of the Status of a Deployment



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

> Deployment SetCurrentDeployment(ctx, setCurrentDeployment)

Set the current Deployment

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**setCurrentDeployment** | [**SetCurrentDeployment**](SetCurrentDeployment.md)|  | 

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


## UpdateDeploymentStatus

> Deployment UpdateDeploymentStatus(ctx, id, deploymentStatusUpdate)

Update parts of the Status of a Deployment

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **int64**|  | 
**deploymentStatusUpdate** | [**DeploymentStatusUpdate**](DeploymentStatusUpdate.md)|  | 

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

