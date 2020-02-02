part of openapi.api;

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
    id = json['id'];
    name = json['name'];
    created = (json['created'] == null) ?
      null :
      DateTime.parse(json['created']);
    steps = (json['steps'] == null) ?
      null :
      SignalcdStep.listFromJson(json['steps']);
    checks = (json['checks'] == null) ?
      null :
      SignalcdCheck.listFromJson(json['checks']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (id != null)
      json['id'] = id;
    if (name != null)
      json['name'] = name;
    if (created != null)
      json['created'] = created == null ? null : created.toUtc().toIso8601String();
    if (steps != null)
      json['steps'] = steps;
    if (checks != null)
      json['checks'] = checks;
    return json;
  }

  static List<SignalcdPipeline> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdPipeline>() : json.map((value) => SignalcdPipeline.fromJson(value)).toList();
  }

  static Map<String, SignalcdPipeline> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdPipeline>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdPipeline.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdPipeline-objects as value to a dart map
  static Map<String, List<SignalcdPipeline>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdPipeline>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdPipeline.listFromJson(value);
       });
     }
     return map;
  }
}

