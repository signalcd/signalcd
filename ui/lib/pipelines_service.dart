import 'package:angular/angular.dart';
import 'package:ui/api.dart';
import 'package:ui/src/api/api.dart';

@Injectable()
class PipelinesService {
  final API _api;

  PipelinesService(this._api);

  Future<List<Pipeline>> pipelines() async {
    return await _api.pipelines.pipelines();
  }

  Future<dynamic> deploy(String id) async {
//    return await _api.pipelines.updateCurrentPipeline(id);
  }
}
