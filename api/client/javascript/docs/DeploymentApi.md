# SignalCd.DeploymentApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getCurrentDeployment**](DeploymentApi.md#getCurrentDeployment) | **GET** /deployments/current | Get the current Deployment
[**listDeployments**](DeploymentApi.md#listDeployments) | **GET** /deployments | List Deployments
[**setCurrentDeployment**](DeploymentApi.md#setCurrentDeployment) | **POST** /deployments/current | Set the current Deployment
[**updateDeploymentStatus**](DeploymentApi.md#updateDeploymentStatus) | **PATCH** /deployments/{id}/status | Update parts of the Status of a Deployment



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


## setCurrentDeployment

> Deployment setCurrentDeployment(setCurrentDeployment)

Set the current Deployment

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.DeploymentApi();
let setCurrentDeployment = new SignalCd.SetCurrentDeployment(); // SetCurrentDeployment | 
apiInstance.setCurrentDeployment(setCurrentDeployment).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **setCurrentDeployment** | [**SetCurrentDeployment**](SetCurrentDeployment.md)|  | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## updateDeploymentStatus

> Deployment updateDeploymentStatus(id, deploymentStatusUpdate)

Update parts of the Status of a Deployment

### Example

```javascript
import SignalCd from 'signal_cd';

let apiInstance = new SignalCd.DeploymentApi();
let id = 789; // Number | 
let deploymentStatusUpdate = new SignalCd.DeploymentStatusUpdate(); // DeploymentStatusUpdate | 
apiInstance.updateDeploymentStatus(id, deploymentStatusUpdate).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **Number**|  | 
 **deploymentStatusUpdate** | [**DeploymentStatusUpdate**](DeploymentStatusUpdate.md)|  | 

### Return type

[**Deployment**](Deployment.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

