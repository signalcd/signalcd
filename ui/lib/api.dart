import 'dart:html';

import 'package:ui/src/api/api.dart';

class API {
  UIServiceApi ui;

  API() {
    String base = window.location.toString();
    if (base.endsWith("/")) {
      base = base.substring(0, base.length - 1);
    }

    ApiClient client = ApiClient(basePath: base);
    this.ui = UIServiceApi(client);
  }
}
