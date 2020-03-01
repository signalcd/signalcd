part of openapi.api;

class SignalcdStatus {
  
  String exitCode = null;
  
  DateTime started = null;
  
  DateTime stopped = null;
  SignalcdStatus();

  @override
  String toString() {
    return 'SignalcdStatus[exitCode=$exitCode, started=$started, stopped=$stopped, ]';
  }

  SignalcdStatus.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    exitCode = json['exitCode'];
    started = (json['started'] == null) ?
      null :
      DateTime.parse(json['started']);
    stopped = (json['stopped'] == null) ?
      null :
      DateTime.parse(json['stopped']);
  }

  Map<String, dynamic> toJson() {
    Map <String, dynamic> json = {};
    if (exitCode != null)
      json['exitCode'] = exitCode;
    if (started != null)
      json['started'] = started == null ? null : started.toUtc().toIso8601String();
    if (stopped != null)
      json['stopped'] = stopped == null ? null : stopped.toUtc().toIso8601String();
    return json;
  }

  static List<SignalcdStatus> listFromJson(List<dynamic> json) {
    return json == null ? List<SignalcdStatus>() : json.map((value) => SignalcdStatus.fromJson(value)).toList();
  }

  static Map<String, SignalcdStatus> mapFromJson(Map<String, dynamic> json) {
    var map = Map<String, SignalcdStatus>();
    if (json != null && json.isNotEmpty) {
      json.forEach((String key, dynamic value) => map[key] = SignalcdStatus.fromJson(value));
    }
    return map;
  }

  // maps a json object with a list of SignalcdStatus-objects as value to a dart map
  static Map<String, List<SignalcdStatus>> mapListFromJson(Map<String, dynamic> json) {
    var map = Map<String, List<SignalcdStatus>>();
     if (json != null && json.isNotEmpty) {
       json.forEach((String key, dynamic value) {
         map[key] = SignalcdStatus.listFromJson(value);
       });
     }
     return map;
  }
}

