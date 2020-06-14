# SignalCd.DeploymentApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCurrentDeployment**](DeploymentApi.md#getCurrentDeployment) | **GET** /deployments/current | Get the current Deployment
[**listDeployments**](DeploymentApi.md#listDeployments) | **GET** /deployments | List Deployments



## getCurrentDeployment

> Deployment getCurrentDeployment()

Get the current Deployment

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.DeploymentApi();
apiInstance.getCurrentDeployment().then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## listDeployments

> [Deployment] listDeployments()

List Deployments

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.DeploymentApi();
apiInstance.listDeployments().then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters

This endpoint does not need any parameter.

### Return type

[**[Deployment]**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

