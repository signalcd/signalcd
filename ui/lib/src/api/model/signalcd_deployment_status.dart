part of openapi.api;

class SignalcdDeploymentStatus {
  
  DeploymentStatusPhase phase = null;
  //enum phaseEnum {  DONTUSETHISVALUE,  UNKNOWN,  SUCCESS,  FAILURE,  PROGRESS,  PENDING,  KILLED,  };{
  SignalcdDeploymentStatus();

  @override
  String toString() {
    return 'SignalcdDeploymentStatus[phase=$phase, ]';
  }

  SignalcdDeploymentStatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    phase = (json['phase'] == null) ?
      null :
      DeploymentStatusPhase.fromJson(json['phase']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (phase != null)
      json['phase'] = phase;
    return json;
  }

  static List<SignalcdDeploymentStatus> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdDeploymentStatus>() : json.map((value) => SignalcdDeploymentStatus.fromJson(value)).toList();
  }

  static Map<String, SignalcdDeploymentStatus> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdDeploymentStatus>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdDeploymentStatus.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdDeploymentStatus-objects as value to a dart map
  static Map<String, List<SignalcdDeploymentStatus>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdDeploymentStatus>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdDeploymentStatus.listFromJson(value);
       });
     }
     return map;
  }
}

