import 'package:ui/src/api/api.dart';

class API {
  UIServiceApi ui;

  API() {
    ApiClient client = ApiClient(basePath: "http://localhost:6060");
    this.ui = UIServiceApi(client);
  }
}
