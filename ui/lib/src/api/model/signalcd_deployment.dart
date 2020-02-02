part of openapi.api;

class SignalcdDeployment {
  
  String number = null;
  
  DateTime created = null;
  
  DateTime started = null;
  
  DateTime finished = null;
  
  SignalcdDeploymentStatus status = null;
  
  SignalcdPipeline pipeline = null;
  SignalcdDeployment();

  @override
  String toString() {
    return 'SignalcdDeployment[number=$number, created=$created, started=$started, finished=$finished, status=$status, pipeline=$pipeline, ]';
  }

  SignalcdDeployment.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    number = json['number'];
    created = (json['created'] == null) ?
      null :
      DateTime.parse(json['created']);
    started = (json['started'] == null) ?
      null :
      DateTime.parse(json['started']);
    finished = (json['finished'] == null) ?
      null :
      DateTime.parse(json['finished']);
    status = (json['status'] == null) ?
      null :
      SignalcdDeploymentStatus.fromJson(json['status']);
    pipeline = (json['pipeline'] == null) ?
      null :
      SignalcdPipeline.fromJson(json['pipeline']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (number != null)
      json['number'] = number;
    if (created != null)
      json['created'] = created == null ? null : created.toUtc().toIso8601String();
    if (started != null)
      json['started'] = started == null ? null : started.toUtc().toIso8601String();
    if (finished != null)
      json['finished'] = finished == null ? null : finished.toUtc().toIso8601String();
    if (status != null)
      json['status'] = status;
    if (pipeline != null)
      json['pipeline'] = pipeline;
    return json;
  }

  static List<SignalcdDeployment> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdDeployment>() : json.map((value) => SignalcdDeployment.fromJson(value)).toList();
  }

  static Map<String, SignalcdDeployment> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdDeployment>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdDeployment.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdDeployment-objects as value to a dart map
  static Map<String, List<SignalcdDeployment>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdDeployment>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdDeployment.listFromJson(value);
       });
     }
     return map;
  }
}

