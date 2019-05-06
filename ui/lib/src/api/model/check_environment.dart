part of swagger.api;

class CheckEnvironment {
  
  String key = null;
  

  String value = null;
  
  CheckEnvironment();

  @override
  String toString() {
    return 'CheckEnvironment[key=$key, value=$value, ]';
  }

  CheckEnvironment.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    key =
        json['key']
    ;
    value =
        json['value']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'key': key,
      'value': value
     };
  }

  static List<CheckEnvironment> listFromJson(List<dynamic> json) {
    return json == null ? new List<CheckEnvironment>() : json.map((value) => new CheckEnvironment.fromJson(value)).toList();
  }

  static Map<String, CheckEnvironment> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, CheckEnvironment>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new CheckEnvironment.fromJson(value));
    }
    return map;
  }
}

