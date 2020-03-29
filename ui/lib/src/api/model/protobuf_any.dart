part of openapi.api;

class ProtobufAny {
  
  String typeUrl = null;
  
  String value = null;
  ProtobufAny();

  @override
  String toString() {
    return 'ProtobufAny[typeUrl=$typeUrl, value=$value, ]';
  }

  ProtobufAny.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    typeUrl = json['type_url'];
    value = json['value'];
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (typeUrl != null)
      json['type_url'] = typeUrl;
    if (value != null)
      json['value'] = value;
    return json;
  }

  static List<ProtobufAny> listFromJson(List<dynamic> json) {
    return json == null ? List<ProtobufAny>() : json.map((value) => ProtobufAny.fromJson(value)).toList();
  }

  static Map<String, ProtobufAny> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, ProtobufAny>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = ProtobufAny.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of ProtobufAny-objects as value to a dart map
  static Map<String, List<ProtobufAny>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<ProtobufAny>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = ProtobufAny.listFromJson(value);
       });
     }
     return map;
  }
}

