part of swagger.api;

class SignalcdCheck {
  
  String name = null;
  

  String image = null;
  

  List<String> imagePullSecrets = [];
  

  String duration = null;
  
  SignalcdCheck();

  @override
  String toString() {
    return 'SignalcdCheck[name=$name, image=$image, imagePullSecrets=$imagePullSecrets, duration=$duration, ]';
  }

  SignalcdCheck.fromJson(Map<String, dynamic> json) {
    if (json == null) return;
    name =
        json['name']
    ;
    image =
        json['image']
    ;
    imagePullSecrets =
        (json['ImagePullSecrets'] as List).map((item) => item as String).toList()
    ;
    duration =
        json['duration']
    ;
  }

  Map<String, dynamic> toJson() {
    return {
      'name': name,
      'image': image,
      'ImagePullSecrets': imagePullSecrets,
      'duration': duration
     };
  }

  static List<SignalcdCheck> listFromJson(List<dynamic> json) {
    return json == null ? new List<SignalcdCheck>() : json.map((value) => new SignalcdCheck.fromJson(value)).toList();
  }

  static Map<String, SignalcdCheck> mapFromJson(Map<String, Map<String, dynamic>> json) {
    var map = new Map<String, SignalcdCheck>();
    if (json != null && json.length > 0) {
      json.forEach((String key, Map<String, dynamic> value) => map[key] = new SignalcdCheck.fromJson(value));
    }
    return map;
  }
}

