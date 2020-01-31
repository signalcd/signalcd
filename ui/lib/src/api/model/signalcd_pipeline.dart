part of swagger.api;

class SignalcdPipeline {
  
  String id = null;
  

  String name = null;
  

  DateTime created = null;
  

  List<SignalcdStep> steps = [];
  

  List<SignalcdCheck> checks = [];
  
  SignalcdPipeline();

  @override
  String toString() {
    return 'SignalcdPipeline[id=$id, name=$name, created=$created, steps=$steps, checks=$checks, ]';
  }

  SignalcdPipeline.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    id =
        json['id']
    ;
    name =
        json['name']
    ;
    created = json['created'] == null ? null : DateTime.parse(json['created']);
    steps =
      SignalcdStep.listFromJson(json['steps'])
;
    checks =
      SignalcdCheck.listFromJson(json['checks'])
;
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'created': created == null ? '' : created.toUtc().toIso8601String(),
      'steps': steps,
      'checks': checks
     };
  }

  static List<SignalcdPipeline> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdPipeline>() : json.map((value) => new SignalcdPipeline.fromJson(value)).toList();
  }

  static Map<String, SignalcdPipeline> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdPipeline>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdPipeline.fromJson(value));
    }
    return map;
  }
}

