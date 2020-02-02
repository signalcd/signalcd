part of openapi.api;

class SignalcdSetCurrentDeploymentResponse {
  
  SignalcdDeployment deployment = null;
  SignalcdSetCurrentDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdSetCurrentDeploymentResponse[deployment=$deployment, ]';
  }

  SignalcdSetCurrentDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployment = (json['deployment'] == null) ?
      null :
      SignalcdDeployment.fromJson(json['deployment']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (deployment != null)
      json['deployment'] = deployment;
    return json;
  }

  static List<SignalcdSetCurrentDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdSetCurrentDeploymentResponse>() : json.map((value) => SignalcdSetCurrentDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdSetCurrentDeploymentResponse> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdSetCurrentDeploymentResponse>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdSetCurrentDeploymentResponse.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdSetCurrentDeploymentResponse-objects as value to a dart map
  static Map<String, List<SignalcdSetCurrentDeploymentResponse>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdSetCurrentDeploymentResponse>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdSetCurrentDeploymentResponse.listFromJson(value);
       });
     }
     return map;
  }
}

