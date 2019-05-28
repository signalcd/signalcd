import 'package:angular/angular.dart';

@Component(
  selector: 'status-icon',
  templateUrl: 'status_component.html',
  directives: [coreDirectives],
)
class StatusComponent {
  @Input()
  String status;
}
