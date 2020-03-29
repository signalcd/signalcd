part of openapi.api;

class RuntimeError {
  
  String error = null;
  
  int code = null;
  
  String message = null;
  
  List<ProtobufAny> details = [];
  RuntimeError();

  @override
  String toString() {
    return 'RuntimeError[error=$error, code=$code, message=$message, details=$details, ]';
  }

  RuntimeError.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    error = json['error'];
    code = json['code'];
    message = json['message'];
    details = (json['details'] == null) ?
      null :
      ProtobufAny.listFromJson(json['details']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (error != null)
      json['error'] = error;
    if (code != null)
      json['code'] = code;
    if (message != null)
      json['message'] = message;
    if (details != null)
      json['details'] = details;
    return json;
  }

  static List<RuntimeError> listFromJson(List<dynamic> json) {
    return json == null ? List<RuntimeError>() : json.map((value) => RuntimeError.fromJson(value)).toList();
  }

  static Map<String, RuntimeError> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, RuntimeError>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = RuntimeError.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of RuntimeError-objects as value to a dart map
  static Map<String, List<RuntimeError>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<RuntimeError>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = RuntimeError.listFromJson(value);
       });
     }
     return map;
  }
}

