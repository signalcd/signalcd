# \PipelineApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetPipeline**](PipelineApi.md#GetPipeline) | **Get** /pipelines/{id} | Get Pipeline by its ID
[**ListPipelines**](PipelineApi.md#ListPipelines) | **Get** /pipelines | List of Pipelines.



## GetPipeline

> Pipeline GetPipeline(ctx, id)

Get Pipeline by its ID

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | [**string**](.md)| Pipeline ID (UUID) | 

### Return type

[**Pipeline**](Pipeline.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPipelines

> []Pipeline ListPipelines(ctx, )

List of Pipelines.

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]Pipeline**](Pipeline.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

