part of swagger.api;

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
    number =
        json['number']
    ;
    created = json['created'] == null ? null : DateTime.parse(json['created']);
    started = json['started'] == null ? null : DateTime.parse(json['started']);
    finished = json['finished'] == null ? null : DateTime.parse(json['finished']);
    status =
      
      
      new SignalcdDeploymentStatus.fromJson(json['status'])
;
    pipeline =
      
      
      new SignalcdPipeline.fromJson(json['pipeline'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'number': number,
      'created': created == null ? '' : created.toUtc().toIso8601String(),
      'started': started == null ? '' : started.toUtc().toIso8601String(),
      'finished': finished == null ? '' : finished.toUtc().toIso8601String(),
      'status': status,
      'pipeline': pipeline
     };
  }

  static List<SignalcdDeployment> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdDeployment>() : json.map((value) => new SignalcdDeployment.fromJson(value)).toList();
  }

  static Map<String, SignalcdDeployment> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdDeployment>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdDeployment.fromJson(value));
    }
    return map;
  }
}

