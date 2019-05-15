part of swagger.api;

class Deployment {
  
  num number = null;
  

  DateTime created = null;
  

  Deploymentstatus status = null;
  

  Pipeline pipeline = null;
  
  Deployment();

  @override
  String toString() {
    return 'Deployment[number=$number, created=$created, status=$status, pipeline=$pipeline, ]';
  }

  Deployment.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    number =
        json['number']
    ;
    created = json['created'] == null ? null : DateTime.parse(json['created']);
    status =
      
      
      new Deploymentstatus.fromJson(json['status'])
;
    pipeline =
      
      
      new Pipeline.fromJson(json['pipeline'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'number': number,
      'created': created == null ? '' : created.toUtc().toIso8601String(),
      'status': status,
      'pipeline': pipeline
     };
  }

  static List<Deployment> listFromJson(List<dynamic> json) {
    return json == null ? new List<Deployment>() : json.map((value) => new Deployment.fromJson(value)).toList();
  }

  static Map<String, Deployment> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, Deployment>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new Deployment.fromJson(value));
    }
    return map;
  }
}

