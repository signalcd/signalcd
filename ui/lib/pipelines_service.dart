import 'package:angular/angular.dart';
import 'package:ui/api.dart';
import 'package:ui/src/api/api.dart';

@Injectable()
class PipelinesService {
  final API _api;

  PipelinesService(this._api);

  Future<List<SignalcdPipeline>> pipelines() async {
    return await _api.ui.listPipelines().then((resp) => resp.pipelines);
  }
}
