part of swagger.api;

class SignalcdSetCurrentDeploymentResponse {
  
  SignalcdDeployment deployment = null;
  
  SignalcdSetCurrentDeploymentResponse();

  @override
  String toString() {
    return 'SignalcdSetCurrentDeploymentResponse[deployment=$deployment, ]';
  }

  SignalcdSetCurrentDeploymentResponse.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    deployment =
      
      
      new SignalcdDeployment.fromJson(json['deployment'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'deployment': deployment
     };
  }

  static List<SignalcdSetCurrentDeploymentResponse> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdSetCurrentDeploymentResponse>() : json.map((value) => new SignalcdSetCurrentDeploymentResponse.fromJson(value)).toList();
  }

  static Map<String, SignalcdSetCurrentDeploymentResponse> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdSetCurrentDeploymentResponse>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdSetCurrentDeploymentResponse.fromJson(value));
    }
    return map;
  }
}

