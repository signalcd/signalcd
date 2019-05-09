import 'package:angular/angular.dart';
import 'package:ui/api.dart';
import 'package:ui/app_component.template.dart' as ng;
import 'main.template.dart' as self;

@GenerateInjector([
  ClassProvider(API),
])
final InjectorFactory injector = self.injector$Injector;

void main() {
  runApp(ng.AppComponentNgFactory, createInjector: injector);
}
