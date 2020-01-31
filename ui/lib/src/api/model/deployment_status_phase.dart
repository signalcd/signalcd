part of swagger.api;

class DeploymentStatusPhase {
  /// The underlying value of this enum member.
  String value;

  DeploymentStatusPhase._internal(this.value);

  static DeploymentStatusPhase dONTUSETHISVALUE_ = DeploymentStatusPhase._internal("DONTUSETHISVALUE");
  static DeploymentStatusPhase uNKNOWN_ = DeploymentStatusPhase._internal("UNKNOWN");
  static DeploymentStatusPhase sUCCESS_ = DeploymentStatusPhase._internal("SUCCESS");
  static DeploymentStatusPhase fAILURE_ = DeploymentStatusPhase._internal("FAILURE");
  static DeploymentStatusPhase pROGRESS_ = DeploymentStatusPhase._internal("PROGRESS");
  static DeploymentStatusPhase pENDING_ = DeploymentStatusPhase._internal("PENDING");
  static DeploymentStatusPhase kILLED_ = DeploymentStatusPhase._internal("KILLED");

  DeploymentStatusPhase.fromJson(dynamic data) {
    switch (data) {
          case "DONTUSETHISVALUE": value = data; break;
          case "UNKNOWN": value = data; break;
          case "SUCCESS": value = data; break;
          case "FAILURE": value = data; break;
          case "PROGRESS": value = data; break;
          case "PENDING": value = data; break;
          case "KILLED": value = data; break;
    default: throw('Unknown enum value to decode: $data');
    }
  }

  static dynamic encode(DeploymentStatusPhase data) {
    return data.value;
  }
}

