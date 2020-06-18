# SignalCd.PipelineApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createPipeline**](PipelineApi.md#createPipeline) | **POST** /pipelines | Create a new Pipeline.
[**getPipeline**](PipelineApi.md#getPipeline) | **GET** /pipelines/{id} | Get Pipeline by its ID
[**listPipelines**](PipelineApi.md#listPipelines) | **GET** /pipelines | List of Pipelines.



## createPipeline

> Pipeline createPipeline(pipeline)

Create a new Pipeline.

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.PipelineApi();
let pipeline = new SignalCd.Pipeline(); // Pipeline | 
apiInstance.createPipeline(pipeline).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pipeline** | [**Pipeline**](Pipeline.md)|  | 

### Return type

[**Pipeline**](Pipeline.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## getPipeline

> Pipeline getPipeline(id)

Get Pipeline by its ID

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.PipelineApi();
let id = null; // String | Pipeline ID (UUID)
apiInstance.getPipeline(id).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | [**String**](.md)| Pipeline ID (UUID) | 

### Return type

[**Pipeline**](Pipeline.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listPipelines

> [Pipeline] listPipelines()

List of Pipelines.

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.PipelineApi();
apiInstance.listPipelines().then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters

This endpoint does not need any parameter.

### Return type

[**[Pipeline]**](Pipeline.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

