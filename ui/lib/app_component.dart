import 'package:angular/angular.dart';
import 'package:ui/pipelines_component.dart';

@Component(
  selector: 'my-app',
  templateUrl: 'app_component.html',
  directives: [
    PipelinesComponent,
  ],
)
class AppComponent {
  var name = 'Angular';
}
