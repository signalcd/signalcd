import 'package:ui/src/api/api.dart';

class API {
  UIServiceApi ui;

  API() {
    ApiClient client = ApiClient(basePath: "/");
    this.ui = UIServiceApi(client);
  }
}
