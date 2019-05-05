import 'package:angular/angular.dart';

@Component(
  selector: 'pipelines-list',
  templateUrl: 'pipelines_component.html',
  directives: [
    coreDirectives,
  ]
)
class PipelinesComponent implements OnInit {
  List<Pipeline> pipelines = [];

  @override
  void ngOnInit() {
    pipelines.add(Pipeline(
      id: 'eee4047d-3826-4bf0-a7f1-b0b339521a52',
      name: 'cheese0',
    ));
    pipelines.add(Pipeline(
      id: '6151e283-99b6-4611-bbc4-8aa4d3ddf8fd',
      name: 'cheese1',
    ));
    pipelines.add(Pipeline(
      id: 'a7cae189-400e-4d8c-a982-f0e9a5b4901f',
      name: 'cheese2',
    ));
  }
}

class Pipeline {
  String id;
  String name;

  Pipeline({this.id, this.name});
}
